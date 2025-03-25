package handlers

import (
	"CountriesDashboardService/consts"
	"CountriesDashboardService/firebase"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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

	// Fetch additional data from the REST Countries API
	url := consts.RestCountriesAPI + "/name/" + content.Country
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

	// Parse the API response
	var apiResponse []map[string]interface{}
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		sendErrorResponse(writer, "Error parsing API response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(apiResponse) == 0 {
		sendErrorResponse(writer, "No data found for the specified country", http.StatusNotFound)
		return
	}

	// Extract the relevant data
	result := make(map[string]interface{})
	countryData := apiResponse[0]

	// Extraction methods
	if content.Features.Capital {
		if capital, ok := extractCapital(countryData); ok {
			result["capital"] = capital
		}
	}

	if content.Features.Coordinates {
		if coordinates, ok := extractCoordinates(countryData); ok {
			result["coordinates"] = coordinates
		}
	}

	if content.Features.Population {
		if population, ok := extractPopulation(countryData); ok {
			result["population"] = population
		}
	}

	if content.Features.Area {
		if area, ok := extractArea(countryData); ok {
			result["area"] = area
		}
	}

	log.Printf("Filtered data: %v", result)

	// Send the filtered response
	if err := json.NewEncoder(writer).Encode(result); err != nil {
		sendErrorResponse(writer, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Successfully retrieved and filtered the dashboard configuration.")
}

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
