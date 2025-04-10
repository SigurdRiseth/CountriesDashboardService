package handlers

import (
	"CountriesDashboardService/cache"
	"CountriesDashboardService/consts"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var GetDashboardConfigFromDBFunc = getDashboardConfigFromDB // Functions to be mocked for testing
var FetchCountryDataFunc = fetchCountryData
var FetchWeatherDataFunc = fetchWeatherData
var FetchCurrencyDataFunc = fetchCurrencyData

// ViewDashboard handles HTTP GET requests to retrieve a populated dashboard configuration from Firestore.
// It fetches the dashboard configuration for the given ID, populates the requested features (e.g., country data,
// weather data, currency exchange rates), and returns the result as a JSON response.
//
// Parameters:
//   - writer: The http.ResponseWriter to write the response to.
//   - request: The http.Request containing the request details, including the "id" query parameter.
//
// The function expects the "id" query parameter to identify the dashboard configuration in Firestore.
// If the ID is missing, or if there are errors fetching data from Firestore or external APIs, it returns an
// appropriate HTTP error status (e.g., 400 Bad Request, 500 Internal Server Error). On success, it returns a
// 200 OK status with the populated dashboard data in JSON format.
//
// Example response:
//
//	  {
//	    "country": "Norway",
//	    "isoCode": "NO",
//	    "features": {
//	      "precipitation": -1.2,
//	      "temperature": 0.80,
//	      "capital": "Oslo",
//	      "coordinates": {
//	        "latitude": 62.0,
//	        "longitude": 10.0
//	      },
//	      "population": 5379475,
//	      "area": 323892.0,
//	      "targetCurrencies": {
//	        "EUR": 0.87701435,
//	        "USD": 0.95184741,
//	        "SEK": 0.9782275
//	      }
//	    },
//	    "lastRetrieval": "20250325 14:00"
//	}
func ViewDashboard(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set(consts.ContentTypeHeader, consts.ApplicationJSON)

	// Extract and validate ID
	id := request.URL.Query().Get(consts.QueryParamID)
	if id == consts.Bunny {
		http.Error(writer, consts.MissingIDParamInURL, http.StatusBadRequest)
		return
	}

	// Fetch configuration from Firestore
	config, err := GetDashboardConfigFromDBFunc(id)
	if err != nil {
		http.Error(writer, consts.FailedRetrieveDBConfig, http.StatusInternalServerError)
		log.Println(consts.LogErrorRetrievingConfig, err)
		return
	}

	// Initialize response
	response := map[string]interface{}{
		consts.CountryString:       config.Country,
		consts.ISOCodeString:       config.IsoCode,
		consts.FeaturesString:      make(map[string]interface{}),
		consts.LastRetrievedString: time.Now().Format(time.RFC3339),
	}
	features := response[consts.FeaturesString].(map[string]interface{})

	// Fetch country data
	countryData := FetchCountryDataFunc(config.Country)
	if countryData == nil {
		log.Println(consts.LogErrorRetrievingCountryData)
		http.Error(writer, consts.FailedToFetchCountryData, http.StatusInternalServerError)
		return
	}

	// Process country-related features
	populateCountryFeatures(features, config, countryData)

	// Process weather-related features
	populateWeatherFeatures(features, config, countryData)

	// Process currency-related features
	if config.Features.TargetCurrencies != nil && len(*config.Features.TargetCurrencies) > 0 {
		fetchCurrencyData(features, countryData, *config.Features.TargetCurrencies)
	}

	// Encode and send response
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		http.Error(writer, consts.FailedEncodeJSON, http.StatusInternalServerError)
		log.Println(consts.ErrorEncodingResponse, err)
	}

	go CheckWebhooks(config.IsoCode, Invoke, id)
	log.Println(consts.LogDBResponseSent, response)
}

// populateCountryFeatures populates country-related features in the dashboard response.
// It extracts data such as the capital, coordinates, population, and area from the country data,
// based on the configuration settings, and adds them to the features map.
//
// Parameters:
//   - features: The map to populate with country-related features (e.g., "capital", "coordinates").
//   - config: The RegistrationRequestBody containing the configuration settings for which features to include.
//   - countryData: The country data fetched from an external API (e.g., REST Countries API), containing fields like
//     "capital", "latlng", "population", and "area".
//
// This function modifies the features map in place. If a feature is not enabled in the config or if the data
// cannot be extracted, the feature is silently skipped.
// populateCountryFeatures populates country-related features (capital, coordinates, population, area) into the features map.
func populateCountryFeatures(features map[string]interface{}, config *consts.RegistrationRequestBody, countryData map[string]interface{}) {
	type featureConfig struct {
		enabled *bool
		extract func(map[string]interface{}) (interface{}, bool)
		key     string
	}

	// Defining the features to process
	featureList := []featureConfig{
		{enabled: config.Features.Capital, extract: func(data map[string]interface{}) (interface{}, bool) {
			return extractCapital(data)
		}, key: consts.CapitalString},
		{enabled: config.Features.Coordinates, extract: func(data map[string]interface{}) (interface{}, bool) {
			return extractCoordinates(data)
		}, key: consts.CoordinatesString},
		{enabled: config.Features.Population, extract: func(data map[string]interface{}) (interface{}, bool) {
			return extractPopulation(data)
		}, key: consts.PopulationString},
		{enabled: config.Features.Area, extract: func(data map[string]interface{}) (interface{}, bool) {
			return extractArea(data)
		}, key: consts.AreaString},
	}
	// Processing each feature
	for _, f := range featureList {
		if f.enabled != nil && *f.enabled {
			if value, ok := f.extract(countryData); ok {
				features[f.key] = value
			}
		}
	}
}

// populateWeatherFeatures populates weather-related features in the dashboard response.
// It fetches weather data (temperature and precipitation) for the country's coordinates if the corresponding features
// are enabled in the configuration.
//
// Parameters:
//   - features: The map to populate with weather-related features (e.g., "temperature", "precipitation").
//   - config: The RegistrationRequestBody containing the configuration settings for which weather features to include.
//   - countryData: The country data fetched from an external API, used to extract coordinates for the weather API request.
//
// This function modifies the features map in place. It first checks if weather data is needed (i.e., if temperature or precipitation is enabled).
// It then extracts the country's coordinates and fetches weather data from the Open-Meteo API if coordinates are available.
// If weather data cannot be fetched or coordinates are unavailable, the weather features are silently skipped.
func populateWeatherFeatures(features map[string]interface{}, config *consts.RegistrationRequestBody, countryData map[string]interface{}) {
	// Check if weather data is needed
	needWeather := (config.Features.Temperature != nil && *config.Features.Temperature) ||
		(config.Features.Precipitation != nil && *config.Features.Precipitation)
	if !needWeather {
		return
	}

	// Extract coordinates for weather data
	var latitude, longitude float64
	if coordinates, ok := extractCoordinates(countryData); ok {
		latitude = coordinates[consts.LatitudeString].(float64)
		longitude = coordinates[consts.LongitudeString].(float64)
	}

	// Fetch weather data if coordinates are available
	if latitude != 0 && longitude != 0 {
		fetchWeatherData(features, latitude, longitude, config)
	}
}

// fetchCountryData retrieves country data from the REST Countries API for the specified country.
// It makes an HTTP GET request to the API and returns the first matching country's data as a map.
//
// Parameters:
//   - countryName: The name of the country to fetch data for (e.g., "Norway").
//
// Returns:
//   - A map containing the country data (e.g., "capital", "latlng", "population", "currencies"), or nil if the request fails or no data is found.
//
// Errors are logged if the HTTP request fails or if the response cannot be decoded. The function returns nil in such cases.
func fetchCountryData(countryName string) map[string]interface{} {
	const maxAge = 6 * time.Hour
	ctx := context.Background()

	// Trying to get cached country info
	if cached, found, err := cache.GetCachedCountryInfo(ctx, countryName, maxAge); err != nil && found {
		log.Println(consts.LogCountryDataCacheHIT, countryName)
		return cached.(map[string]interface{})
	} else if err != nil {
		log.Println(consts.LogCountryDataCacheMISS, countryName, "-", err)
	}

	// Falling to fetching from external API
	url := consts.RestCountriesAPI + consts.QueryParamName + countryName
	resp, err := http.Get(url)
	if err != nil {
		log.Println(consts.FailedToFetchCountryData, err)
		return nil
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(consts.ClosingResponseBody, err)
		}
	}()

	var data []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil || len(data) == 0 {
		log.Println(consts.LogErrorDecodeCurrency, err)
		return nil
	}

	countryData := data[0]

	// Saving to cache
	if err := cache.SaveCountryInfoToCache(countryName, countryData); err != nil {
		log.Println(consts.LogErrorCacheCountryData, err)
	}

	return countryData
}

// fetchWeatherData retrieves weather data (temperature and precipitation) from the Open-Meteo API.
// It uses the provided latitude and longitude to fetch hourly weather data and calculates the average temperature and
// precipitation if enabled in the configuration.
//
// Parameters:
//   - features: The map to populate with weather data (e.g., "temperature", "precipitation").
//   - lat: The latitude of the location to fetch weather data for.
//   - lon: The longitude of the location to fetch weather data for.
//   - config: The RegistrationRequestBody containing the configuration settings for which weather features to include.
//
// This function modifies the features map in place. It fetches hourly weather data and calculates averages for
// temperature and precipitation.
// If the request fails, the response cannot be decoded, or the expected data is not found,
// the function logs an error and returns without modifying the features map.
func fetchWeatherData(features map[string]interface{}, lat, lon float64, config *consts.RegistrationRequestBody) {
	const maxAge = 6 * time.Hour
	ctx := context.Background()
	cacheKey := cache.WeatherCacheKey(lat, lon)

	// Try cache
	if data := getCachedWeatherData(ctx, cacheKey, maxAge); data != nil {
		log.Println(consts.WeatherCacheHit, cacheKey)
		populateWeatherFromMap(features, data, config)
		return
	}

	// Fallback to API
	log.Println(consts.WeatherCacheMiss, cacheKey)
	data := fetchWeatherFromAPI(lat, lon)
	if data == nil {
		log.Println(consts.LogFailedFetchWeatherAPI)
		return
	}

	// Cache it
	if err := cache.SaveWeatherToCache(cacheKey, data); err != nil {
		log.Println(consts.FailCacheWeatherData, err)
	}

	// Populate response
	populateWeatherFromMap(features, data, config)
}

// getCachedWeatherData retrieves cached weather data from the cache.
func getCachedWeatherData(ctx context.Context, key string, maxAge time.Duration) map[string]interface{} {
	raw, found, err := cache.GetCachedWeather(ctx, key, maxAge)
	if err != nil || !found {
		return nil
	}
	if casted, ok := raw.(map[string]interface{}); ok {
		return casted
	}
	return nil
}

// fetchWeatherFromAPI retrieves weather data from the Open-Meteo API.
func fetchWeatherFromAPI(lat, lon float64) map[string]interface{} {
	url := fmt.Sprintf(consts.WeatherURL, consts.OpenMeteoAPI, lat, lon)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(consts.LogHTTPReqFAil, err)
		return nil
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(consts.ClosingResponseBody, err)
		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	var parsed map[string]interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		log.Println(consts.LogFailedJSONDecode, err)
		return nil
	}
	return parsed
}

// populateWeatherFromMap populates weather-related features in the features map.
func populateWeatherFromMap(features map[string]interface{}, data map[string]interface{}, config *consts.RegistrationRequestBody) {
	hourly, ok := data[consts.HourlyString].(map[string]interface{})
	if !ok {
		log.Println(consts.LogHourlyDataMissingWeather)
		return
	}

	if temps, ok := hourly[consts.TemperatureString2].([]interface{}); ok && *config.Features.Temperature {
		features[consts.TemperatureString] = calculateAverage(temps)
	}
	if precs, ok := hourly[consts.PrecipitationString].([]interface{}); ok && *config.Features.Precipitation {
		features[consts.PrecipitationString] = calculateAverage(precs)
	}
}

// fetchCurrencyData retrieves currency exchange rates from the Currency API.
// It extracts the country's currency code, fetches exchange rates for that currency, and filters the rates based
// on the target currencies specified in the configuration.
//
// Parameters:
//   - features: The map to populate with currency data (e.g., "targetCurrencies").
//   - countryData: The country data fetched from an external API, used to extract the currency code.
//   - targetCurrencies: The list of target currency codes to include in the response (e.g., ["EUR", "USD"]).
//
// This function modifies the features map in place by adding a "targetCurrencies" field with the filtered exchange rates.
// If the currency code cannot be extracted, the request fails, or the response cannot be decoded,
// the function logs an error and returns without modifying the features map.
// The function also logs the API URL and raw response for debugging purposes.
func fetchCurrencyData(features map[string]interface{}, countryData map[string]interface{}, targetCurrencies []string) {
	const maxAge = 6 * time.Hour
	ctx := context.Background()

	currencyCode := extractCurrencyCode(countryData)
	if currencyCode == consts.Bunny {
		log.Println(consts.LogNoCurrencyCodeForCountry)
		return
	}

	// 1. Try cache first
	if tryCurrencyFromCache(features, ctx, currencyCode, targetCurrencies, maxAge) {
		return
	}

	// 2. Fallback to API
	data := fetchCurrencyFromAPI(currencyCode)
	if data == nil {
		return
	}

	// 3. Extract + Save + Populate
	if rates := extractRates(data); rates != nil {
		err := cache.SaveCurrencyRatesToCache(currencyCode, rates)
		if err != nil {
			return
		}
		populateCurrencyFeatures(features, rates, targetCurrencies)
	}
}

// tryCurrencyFromCache attempts to retrieve cached currency rates for the specified currency code.
func tryCurrencyFromCache(features map[string]interface{}, ctx context.Context, currencyCode string, targetCurrencies []string, maxAge time.Duration) bool {
	cached, found, err := cache.GetCachedCurrencyRates(ctx, currencyCode, maxAge)
	if err != nil || !found {
		log.Println(consts.LogCurrencyCacheMiss, currencyCode)
		return false
	}

	typedRates, ok := cached.(map[string]float64)
	if !ok {
		log.Println(consts.LogInvalidCacheTypesCurrency)
		return false
	}

	log.Println(consts.LogCurrencyCacheHIT, currencyCode)
	populateCurrencyFeatures(features, typedRates, targetCurrencies)
	return true
}

// fetchCurrencyFromAPI retrieves currency exchange rates from the Currency API for the specified currency code.
func fetchCurrencyFromAPI(currencyCode string) map[string]interface{} {
	url := fmt.Sprintf(consts.SS, consts.CurrencyAPI, currencyCode)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(consts.LogErrorFetchCurrency, err)
		return nil
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(consts.ClosingResponseBody, err)
		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	var parsed map[string]interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		log.Println(consts.LogErrorDecodeCurrency, err)
		return nil
	}
	return parsed
}

// extractRates extracts exchange rates from the Currency API response.
func extractRates(data map[string]interface{}) map[string]float64 {
	rawRates, ok := data[consts.RatesString].(map[string]interface{})
	if !ok {
		log.Println(consts.LogMissingRatesCurrencyAPI)
		return nil
	}

	rates := make(map[string]float64)
	for k, v := range rawRates {
		if floatVal, ok := v.(float64); ok {
			rates[k] = floatVal
		}
	}
	return rates
}

// populateCurrencyFeatures populates the "targetCurrencies" field in the features map.
func populateCurrencyFeatures(features map[string]interface{}, rates map[string]float64, targetCurrencies []string) {
	filtered := make(map[string]interface{})
	for _, cur := range targetCurrencies {
		if val, exists := rates[cur]; exists {
			filtered[cur] = val
		}
	}
	features[consts.TargetCurrenciesString] = filtered
}

// extractCurrencyCode extracts the ISO 4217 currency code from country data.
// It looks for the "currencies" field in the country data and returns the first currency code found.
//
// Parameters:
//   - countryData: The country data map containing a "currencies" field.
//
// Returns:
//   - The ISO 4217 currency code as a string (e.g., "NOK"), or an empty string if no currency code is found.
func extractCurrencyCode(countryData map[string]interface{}) string {
	if currencies, ok := countryData[consts.CurrenciesString].(map[string]interface{}); ok {
		for code := range currencies {
			return code // Return the first currency code found
		}
	}
	return ""
}

// extractCapital extracts the capital city from the country data.
func extractCapital(countryData map[string]interface{}) (string, bool) {
	if capitalList, ok := countryData[consts.CapitalString].([]interface{}); ok && len(capitalList) > 0 {
		if capital, ok := capitalList[0].(string); ok {
			return capital, true
		}
	}
	return "", false
}

// extractCoordinates extracts latitude and longitude from the country data.
func extractCoordinates(countryData map[string]interface{}) (map[string]interface{}, bool) {
	if latlng, ok := countryData[consts.LatlngString].([]interface{}); ok && len(latlng) == 2 {
		return map[string]interface{}{
			consts.LatitudeString:  latlng[0],
			consts.LongitudeString: latlng[1],
		}, true
	}
	return nil, false
}

// extractPopulation extracts population from the country data.
func extractPopulation(countryData map[string]interface{}) (float64, bool) {
	if population, ok := countryData[consts.PopulationString].(float64); ok {
		return population, true
	}
	return 0, false
}

// extractArea extracts area from the country data.
func extractArea(countryData map[string]interface{}) (float64, bool) {
	if area, ok := countryData[consts.AreaString].(float64); ok {
		return area, true
	}
	return 0, false
}

// Calculate the average of a slice of float64 values
func calculateAverage(data []interface{}) float64 {
	sum := 0.0
	for _, v := range data {
		if value, ok := v.(float64); ok {
			sum += value
		}
	}
	if len(data) == 0 {
		return 0
	}
	return sum / float64(len(data))
}
