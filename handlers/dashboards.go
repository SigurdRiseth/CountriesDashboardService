// TODO: Implement the ViewDashboard function us get method from registrations.
package handlers

import (
	"CountriesDashboardService/consts"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// ViewDashboard handles HTTP GET requests to retrieve a populated dashboard configuration from Firestore.
func ViewDashboard(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	// Extract and validate ID
	id := request.URL.Query().Get("id")
	if id == "" {
		http.Error(writer, "Missing 'id' in URL parameters", http.StatusBadRequest)
		return
	}

	// Fetch configuration from Firestore
	config, err := getDashboardConfigFromDB(id)
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
		"lastRetrieval": time.Now().Format("20060102 15:04"),
	}
	features := response["features"].(map[string]interface{})

	// Fetch country data
	countryData := fetchCountryData(config.Country)
	if countryData == nil {
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

// populateCountryFeatures populates country-related features (capital, coordinates, population, area).
func populateCountryFeatures(features map[string]interface{}, config *consts.RegistrationRequestBody, countryData map[string]interface{}) {
	if config.Features.Capital != nil && *config.Features.Capital {
		if capital, ok := extractCapital(countryData); ok {
			features["capital"] = capital
		}
	}

	if config.Features.Coordinates != nil && *config.Features.Coordinates {
		if coordinates, ok := extractCoordinates(countryData); ok {
			features["coordinates"] = coordinates
		}
	}

	if config.Features.Population != nil && *config.Features.Population {
		if population, ok := extractPopulation(countryData); ok {
			features["population"] = population
		}
	}

	if config.Features.Area != nil && *config.Features.Area {
		if area, ok := extractArea(countryData); ok {
			features["area"] = area
		}
	}
}

// populateWeatherFeatures populates weather-related features (temperature, precipitation).
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

// fetchCountryData makes an API request to the REST Countries API
func fetchCountryData(countryName string) map[string]interface{} {
	url := consts.RestCountriesAPI + "/name/" + countryName
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching country data:", err)
		return nil
	}
	defer resp.Body.Close()

	var data []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil || len(data) == 0 {
		log.Println("Error decoding country data:", err)
		return nil
	}

	return data[0]
}

// fetchWeatherData fetches temperature and precipitation from Open-Meteo API
func fetchWeatherData(features map[string]interface{}, lat, lon float64, config *consts.RegistrationRequestBody) {
	weatherURL := fmt.Sprintf("%s?latitude=%.2f&longitude=%.2f&hourly=temperature_2m,precipitation", consts.OpenMeteoAPI, lat, lon)
	resp, err := http.Get(weatherURL)
	if err != nil {
		log.Println("Error fetching weather data:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
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

// fetchCurrencyData fetches currency exchange rates from Currency API
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
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
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
