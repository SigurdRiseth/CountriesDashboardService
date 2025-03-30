package handlers

import (
	"CountriesDashboardService/consts"
	"CountriesDashboardService/firebase"
	"CountriesDashboardService/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"google.golang.org/api/iterator"
)

// HandleStatus routes requests to the appropriate handler based on the HTTP method.
//
// It handles requests to the /dashboard/v1/status/ endpoint, supporting GET requests to retrieve
// the status of dependent services. Other HTTP methods return a 405 Method Not Allowed response.
//
// Parameters:
//   - writer: The http.ResponseWriter to write the response to.
//   - request: The http.Request containing the request details.
func HandleStatus(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		handleGetStatus(writer, request)
	default:
		sendErrorResponseStatus(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGetStatus handles GET requests to retrieve the status of dependent services.
//
// It checks the availability of dependent services (REST Countries API, Open-Meteo API, Currency API,
// and the Firestore notification database), counts the number of registered webhooks, and calculates
// the service uptime. The response is returned as a JSON object with the following fields:
//   - countries_api: HTTP status code for the REST Countries API.
//   - meteo_api: HTTP status code for the Open-Meteo API.
//   - currency_api: HTTP status code for the Currency API.
//   - notification_db: HTTP status code for the Firestore notification database.
//   - webhooks: Number of registered webhooks in Firestore.
//   - version: Service version ("v1").
//   - uptime: Time in seconds since the service started.
//
// Parameters:
//   - writer: The http.ResponseWriter to write the response to.
//   - request: The http.Request containing the request details.
//
// Returns a 200 OK response with the status data in JSON format, or an appropriate error response
// (e.g., 500 Internal Server Error) if an error occurs while counting webhooks or encoding the response.
//
// Example response:
//
//	{
//	  "countries_api": 404,
//	  "meteo_api": 200,
//	  "currency_api": 200,
//	  "notification_db": 200,
//	  "webhooks": 5,
//	  "version": "v1",
//	  "uptime": 3600
//	}
func handleGetStatus(writer http.ResponseWriter, request *http.Request) {
	log.Println("Processing GET request for service status")

	// Initialize the status response
	status := map[string]interface{}{
		"version": "v1",
		"uptime":  int64(time.Since(utils.StartTime).Seconds()),
	}

	// Check availability of dependent services
	status["countries_api"] = checkServiceAvailability(fmt.Sprintf("%s/name/unknown", consts.RestCountriesAPI))
	status["meteo_api"] = checkServiceAvailability(consts.OpenMeteoAPI + "?latitude=0&longitude=0&hourly=temperature_2m")
	status["currency_api"] = checkServiceAvailability(fmt.Sprintf("%s/NOK", consts.CurrencyAPI))
	status["notification_db"] = checkNotificationDB()

	// Count registered webhooks
	webhookCount, err := countWebhooks()
	if err != nil {
		sendErrorResponse(writer, "Failed to count registered webhooks: "+err.Error(), http.StatusInternalServerError)
		return
	}
	status["webhooks"] = webhookCount

	// Send the response
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(status); err != nil {
		sendErrorResponse(writer, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

// checkServiceAvailability checks the availability of a given external service URL.
//
// It automatically selects the appropriate HTTP method based on the service type:
//   - If the URL appears to be for the Currency API (contains "/currency/"), a GET request is used.
//   - For all other URLs, a lightweight HEAD request is used to minimize data transfer.
//
// The function returns the HTTP status code from the response. If the request fails or
// the service is unreachable, it returns 503 (Service Unavailable).
//
// Parameters:
//   - url: The full URL of the service to check.
//
// Returns:
//   - int: The HTTP status code returned by the service, or 503 if the request fails.
func checkServiceAvailability(url string) int {
	// Use GET if the URL is for the Currency API (based on path pattern)
	isCurrencyAPI := false
	urlContainsCurrency := "/currency/"

	if len(url) > 0 && urlContainsCurrency != "" && utils.Contains(url, urlContainsCurrency) {
		isCurrencyAPI = true
	}

	var resp *http.Response
	var err error

	if isCurrencyAPI {
		resp, err = http.Get(url)
	} else {
		resp, err = http.Head(url)
	}

	if err != nil {
		log.Printf("Error checking service availability for %s: %v", url, err)
		return http.StatusServiceUnavailable
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("Error closing response body for", url, ":", err)
		}
	}()

	return resp.StatusCode
}

// checkNotificationDB checks the availability of the Firestore database by performing a simple query.
//
// It attempts to query the "notifications" collection in Firestore to verify that the database is
// accessible. If the query succeeds (even if no documents are found), it returns a 200 OK status.
// If the Firestore client is not initialized or the query fails, it returns a 503 Service Unavailable status.
//
// Returns:
//   - 200 if the database is accessible, 503 if it is not.
func checkNotificationDB() int {
	if firebase.Client == nil {
		log.Println("Firestore client is not initialized")
		return http.StatusServiceUnavailable // 503 if the client is not initialized
	}

	// Perform a simple query to check database availability
	_, err := firebase.Client.Collection(firebase.NotificationsCollection).Limit(1).Documents(firebase.Ctx).Next()
	if err != nil {
		if err == iterator.Done {
			// No documents found, but the database is accessible
			return http.StatusOK // 200 if the query succeeds (even if no documents exist)
		}
		return http.StatusServiceUnavailable // 503 if the query fails
	}

	return http.StatusOK // 200 if the query succeeds
}

// countWebhooks counts the number of registered webhooks in the Firestore "notifications" collection.
//
// It queries the "notifications" collection in Firestore and returns the total number of documents,
// which represents the number of registered webhooks.
//
// Returns:
//   - The number of webhooks as an int64.
//   - An error if the Firestore client is not initialized or the query fails.
func countWebhooks() (int64, error) {
	if firebase.Client == nil {
		return 0, fmt.Errorf("firestore client is not initialized")
	}

	// Query the notifications collection and count the documents
	docs, err := firebase.Client.Collection(firebase.NotificationsCollection).Documents(firebase.Ctx).GetAll()
	if err != nil {
		return 0, fmt.Errorf("error querying webhooks: %v", err)
	}

	return int64(len(docs)), nil
}

// sendErrorResponse sends an error response with the given message and status code.
//
// It constructs a JSON error response with the specified message and writes it to the response writer
// with the given HTTP status code.
//
// Parameters:
//   - writer: The http.ResponseWriter to write the response to.
//   - message: The error message to include in the response.
//   - statusCode: The HTTP status code to use for the response.
func sendErrorResponseStatus(writer http.ResponseWriter, message string, statusCode int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	response := map[string]string{"error": message}
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Printf("Error encoding error response: %v", err)
	}
}
