package handlers

import (
	"CountriesDashboardService/consts"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Stub-able function vars
var GetDashboardConfigFromDBFunc = getDashboardConfigFromDB

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
	writer.Header().Set("Content-Type", "application/json")

	// Extract and validate ID
	id := request.URL.Query().Get("id")
	if id == "" {
		http.Error(writer, "Missing 'id' in URL parameters", http.StatusBadRequest)
		return
	}

	// Fetch configuration from Firestore
	config, err := GetDashboardConfigFromDBFunc(id)
	if err != nil {
		http.Error(writer, "Failed to retrieve dashboard configuration from database", http.StatusInternalServerError)
		log.Println("Error retrieving configuration:", err)
		return
	}

	// Initialize response
	response := map[string]interface{}{
		"country":       config.Country,
		"isoCode":       config.IsoCode,
		"features":      make(map[string]interface{}),
		"lastRetrieval": time.Now().Format(time.RFC3339),
	}
	features := response["features"].(map[string]interface{})

	// Fetch country data
	countryData := fetchCountryData(config.Country)
	if countryData == nil {
		log.Println("Error retrieving country data")
		http.Error(writer, "Failed to fetch country data", http.StatusInternalServerError)
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
		http.Error(writer, "Failed to encode response as JSON", http.StatusInternalServerError)
		log.Println("Error encoding response:", err)
	}
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
		}, key: "capital"},
		{enabled: config.Features.Coordinates, extract: func(data map[string]interface{}) (interface{}, bool) {
			return extractCoordinates(data)
		}, key: "coordinates"},
		{enabled: config.Features.Population, extract: func(data map[string]interface{}) (interface{}, bool) {
			return extractPopulation(data)
		}, key: "population"},
		{enabled: config.Features.Area, extract: func(data map[string]interface{}) (interface{}, bool) {
			return extractArea(data)
		}, key: "area"},
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
		latitude = coordinates["latitude"].(float64)
		longitude = coordinates["longitude"].(float64)
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
	url := consts.RestCountriesAPI + "/name/" + countryName
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching country data:", err)
		return nil
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("Error closing response body:", err)
		}
	}()

	var data []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil || len(data) == 0 {
		log.Println("Error decoding country data:", err)
		return nil
	}

	return data[0]
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
	weatherURL := fmt.Sprintf("%s?latitude=%.2f&longitude=%.2f&hourly=temperature_2m,precipitation", consts.OpenMeteoAPI, lat, lon)
	resp, err := http.Get(weatherURL)
	if err != nil {
		log.Println("Error fetching weather data:", err)
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("Error closing response body:", err)
		}
	}()

	body, _ := io.ReadAll(resp.Body)
	var weatherData map[string]interface{}
	if err := json.Unmarshal(body, &weatherData); err != nil {
		log.Println("Error decoding weather data:", err)
		return
	}

	hourlyData, ok := weatherData["hourly"].(map[string]interface{})
	if !ok {
		return
	}

	if temps, ok := hourlyData["temperature_2m"].([]interface{}); ok && *config.Features.Temperature {
		features["temperature"] = calculateAverage(temps)
	}
	if precs, ok := hourlyData["precipitation"].([]interface{}); ok && *config.Features.Precipitation {
		features["precipitation"] = calculateAverage(precs)
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
	currencyCode := extractCurrencyCode(countryData)
	if currencyCode == "" {
		log.Println("No valid currency code found for the country")
		return
	}

	url := fmt.Sprintf("%s/%s", consts.CurrencyAPI, currencyCode)
	log.Printf("Currency API URL: %s", url) // Log the URL for debugging

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching currency data:", err)
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("Error closing response body:", err)
		}
	}()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Currency API Raw Response: %s", string(body)) // Log the raw response for debugging

	var currencyData map[string]interface{}
	if err := json.Unmarshal(body, &currencyData); err != nil {
		log.Println("Error decoding currency data:", err)
		return
	}

	if rates, ok := currencyData["rates"].(map[string]interface{}); ok {
		filteredRates := make(map[string]interface{})
		for _, currency := range targetCurrencies {
			if rate, exists := rates[currency]; exists {
				filteredRates[currency] = rate
			}
		}
		features["targetCurrencies"] = filteredRates
	}
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
	if currencies, ok := countryData["currencies"].(map[string]interface{}); ok {
		for code := range currencies {
			return code // Return the first currency code found
		}
	}
	return ""
}

// extractCapital extracts the capital city from the country data.
func extractCapital(countryData map[string]interface{}) (string, bool) {
	if capitalList, ok := countryData["capital"].([]interface{}); ok && len(capitalList) > 0 {
		if capital, ok := capitalList[0].(string); ok {
			return capital, true
		}
	}
	return "", false
}

// extractCoordinates extracts latitude and longitude from the country data.
func extractCoordinates(countryData map[string]interface{}) (map[string]interface{}, bool) {
	if latlng, ok := countryData["latlng"].([]interface{}); ok && len(latlng) == 2 {
		return map[string]interface{}{
			"latitude":  latlng[0],
			"longitude": latlng[1],
		}, true
	}
	return nil, false
}

// extractPopulation extracts population from the country data.
func extractPopulation(countryData map[string]interface{}) (float64, bool) {
	if population, ok := countryData["population"].(float64); ok {
		return population, true
	}
	return 0, false
}

// extractArea extracts area from the country data.
func extractArea(countryData map[string]interface{}) (float64, bool) {
	if area, ok := countryData["area"].(float64); ok {
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
