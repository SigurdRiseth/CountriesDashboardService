package handlers

import (
	"CountriesDashboardService/consts"
	"CountriesDashboardService/firebase"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// ViewDashboard handles HTTP GET requests to retrieve a populated dashboard configuration from Firestore.
//
// This function performs the following steps:
// 1. Extracts the configuration ID from the URL query.
// 2. Validates that an ID is provided.
// 3. Retrieves the document from Firestore using the provided ID.
// 4. Fetches additional data from the REST Countries API.
// 5. Returns the structured JSON response if the retrieval is successful.
//
// Parameters:
//   - writer: `http.ResponseWriter`
//     The HTTP response writer used to send data back to the client.
//   - request: `*http.Request`
//     The incoming HTTP request, which must include a valid configuration ID in the URL query.
//
// Behavior:
//   - If no ID is provided in the URL, returns a `400 Bad Request` status.
//   - If the Firebase client is not initialized, returns a `500 Internal Server Error` status.
//   - If the document doesn't exist, returns a `404 Not Found` status.
//   - If retrieval fails for other reasons, returns a `500 Internal Server Error` status.
//   - Upon successful retrieval, returns a `200 OK` status with the populated dashboard in JSON format.
func ViewDashboard(writer http.ResponseWriter, request *http.Request) {
	log.Println("Retrieving populated dashboard...")
	writer.Header().Set("Content-Type", "application/json")

	// Extract the ID from the URL query
	id := request.URL.Query().Get("id")
	if id == "" {
		sendErrorResponse(writer, "Missing ID in URL query", http.StatusBadRequest)
		return
	}

	// Check if Firebase client is initialized
	if firebase.Client == nil {
		sendErrorResponse(writer, "Internal server error: Database client is unavailable", http.StatusInternalServerError)
		return
	}

	// Retrieve the document from Firestore
	docRef := firebase.Client.Collection(collection).Doc(id)
	doc, err := docRef.Get(firebase.Ctx)
	if err != nil {
		sendErrorResponse(writer, "Error retrieving document: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert Firestore document to struct
	var content consts.RegistrationRequestBody
	if err := doc.DataTo(&content); err != nil {
		sendErrorResponse(writer, "Error decoding document data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare the result object
	result := map[string]interface{}{
		"isoCode":       content.IsoCode,
		"country":       content.Country,
		"features":      map[string]interface{}{},
		"lastRetrieval": time.Now().Format("20060102 15:04"),
	}

	features := result["features"].(map[string]interface{})

	// Fetch data from REST Countries API
	url := fmt.Sprintf("%s/name/%s", consts.RestCountriesAPI, content.Country)
	resp, err := http.Get(url)
	if err != nil {
		sendErrorResponse(writer, "Error fetching data from REST Countries API: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		sendErrorResponse(writer, "Error reading response from REST Countries API: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var apiResponse []map[string]interface{}
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		sendErrorResponse(writer, "Error parsing API response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(apiResponse) == 0 {
		sendErrorResponse(writer, "No data found for the specified country", http.StatusNotFound)
		return
	}

	countryData := apiResponse[0]

	if content.Features.Capital {
		if capital, ok := extractCapital(countryData); ok {
			features["capital"] = capital
		}
	}

	if content.Features.Coordinates {
		if coordinates, ok := extractCoordinates(countryData); ok {
			features["coordinates"] = coordinates
		}
	}

	if content.Features.Population {
		if population, ok := extractPopulation(countryData); ok {
			features["population"] = population
		}
	}

	if content.Features.Area {
		if area, ok := extractArea(countryData); ok {
			features["area"] = area
		}
	}

	if content.Features.TargetCurrencies != nil && len(content.Features.TargetCurrencies) > 0 {
		if currencies, ok := extractCurrencyCodes(countryData); ok {
			features["targetCurrencies"] = currencies
		}
	}

	var latitude, longitude float64
	if content.Features.Coordinates {
		if coordinates, ok := extractCoordinates(countryData); ok {
			features["coordinates"] = coordinates
			latitude = coordinates["latitude"].(float64)
			longitude = coordinates["longitude"].(float64)
		}
	}

	if (content.Features.Temperature || content.Features.Precipitation) && latitude != 0 && longitude != 0 {
		meteoURL := fmt.Sprintf("%s?latitude=%.2f&longitude=%.2f&hourly=temperature_2m,precipitation", consts.OpenMeteoAPI, latitude, longitude)
		meteoResp, err := http.Get(meteoURL)
		if err != nil {
			log.Println("Error fetching data from Open-Meteo API:", err)
		} else {
			defer meteoResp.Body.Close()
			meteoBody, _ := ioutil.ReadAll(meteoResp.Body)

			var meteoData map[string]interface{}
			if err := json.Unmarshal(meteoBody, &meteoData); err == nil {
				if hourlyData, ok := meteoData["hourly"].(map[string]interface{}); ok {
					if temps, ok := hourlyData["temperature_2m"].([]interface{}); ok && content.Features.Temperature {
						temperature := calculateAverage(temps)
						features["temperature"] = temperature
					}
					if precs, ok := hourlyData["precipitation"].([]interface{}); ok && content.Features.Precipitation {
						precipitation := calculateAverage(precs)
						features["precipitation"] = precipitation
					}
				}
			}
		}
	}

	log.Printf("Filtered result: %v", result)

	// Send the filtered response
	if err := json.NewEncoder(writer).Encode(result); err != nil {
		sendErrorResponse(writer, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Successfully retrieved and filtered the dashboard configuration.")
}

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

// Methods to extract specific fields from the REST Countries API response

func extractCapital(countryData map[string]interface{}) (interface{}, bool) {
	if capital, ok := countryData["capital"].([]interface{}); ok && len(capital) > 0 {
		return capital[0], true
	}
	return nil, false
}

func extractCoordinates(countryData map[string]interface{}) (map[string]interface{}, bool) {
	if latlng, ok := countryData["latlng"].([]interface{}); ok && len(latlng) == 2 {
		return map[string]interface{}{
			"latitude":  latlng[0],
			"longitude": latlng[1],
		}, true
	}
	return nil, false
}

func extractPopulation(countryData map[string]interface{}) (float64, bool) {
	if population, ok := countryData["population"].(float64); ok {
		return population, true
	}
	return 0, false
}

func extractArea(countryData map[string]interface{}) (float64, bool) {
	if area, ok := countryData["area"].(float64); ok {
		return area, true
	}
	return 0, false
}

// extractCurrencyCodes extracts the currency codes from the country data
func extractCurrencyCodes(countryData map[string]interface{}) ([]string, bool) {
	if currencies, ok := countryData["currencies"].(map[string]interface{}); ok {
		var currencyCodes []string
		for code := range currencies {
			currencyCodes = append(currencyCodes, code)
		}
		return currencyCodes, true
	}
	return nil, false
}
