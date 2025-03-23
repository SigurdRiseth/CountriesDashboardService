package handlers

import (
	"CountriesDashboardService/consts"
	"CountriesDashboardService/firebase"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Collection name in Firestore
const collection = "registrations"

// HandleRegistrations handles HTTP requests for dashboard configurations.
// It supports the following methods:
// - POST: Adds a new dashboard configuration.
// - GET: Views all dashboard configurations or a specific one if an {id} parameter is provided.
// - HEAD: Checks if a specific dashboard configuration exists.
// - PUT: Replaces the existing dashboard configuration.
// - DELETE: Deletes the current dashboard configuration.
// For unsupported methods, it responds with a 405 Method Not Allowed status.
func HandleRegistrations(writer http.ResponseWriter, request *http.Request) {
	log.Println("Registrations endpoint received " + request.Method + " request.")
	switch request.Method {
	case http.MethodPost:
		addDashboardConfiguration(writer, request)
	case http.MethodGet:
		viewDashboardConfiguration(writer, request)
	case http.MethodHead:
		handleHeadRequest(writer, request)
	case http.MethodPut:
		replaceDashboardConfiguration(writer, request)
	case http.MethodDelete:
		deleteDashboardConfiguration(writer, request)
	default:
		log.Printf("Unsupported request method: %s", request.Method)
		http.Error(writer, "Unsupported request method "+request.Method, http.StatusMethodNotAllowed)
		return
	}
}

// addDashboardConfiguration handles HTTP POST requests to add a new dashboard configuration to Firestore.
//
// This function performs the following steps:
// 1. Sets the "Content-Type" header of the response to "application/json".
// 2. Decodes the incoming JSON payload into a `RegistrationRequestBody` struct.
// 3. Adds the current timestamp to the `TimeChanged` field of the decoded struct.
// 4. Writes the decoded and timestamped struct to the Firestore collection defined by the `collection` constant.
// 5. Returns a JSON response containing the document ID and timestamp if the Firestore operation is successful.
//
// Parameters:
//   - writer: `http.ResponseWriter`
//     The HTTP response writer used to send data back to the client.
//   - request: `*http.Request`
//     The incoming HTTP request, which must contain a valid JSON payload in the request body.
//
// Behavior:
//   - If the request body cannot be decoded into the expected struct, the function returns a `400 Bad Request` status
//     and logs an error.
//   - If the Firestore client is not initialized, it returns a `500 Internal Server Error` status and logs an error.
//   - If Firestore fails to add the document, it returns a `400 Bad Request` status with the corresponding error message.
//   - Upon successful addition, the function returns a `201 Created` status with a JSON response containing:
//   - `id`: The Firestore-generated document ID.
//   - `lastChange`: The timestamp when the document was added.
//
// Example Response:
//
//	If the document is successfully added to Firestore, the response might look like:
//	```json
//	{
//	  "id": "D2ASrlTSyTDeYAz00bLf",
//	  "lastChange": "2025-03-23T17:29:44+01:00"
//	}
//	```
func addDashboardConfiguration(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	// Decode the request body into a RegistrationRequestBody struct
	var content consts.RegistrationRequestBody
	if err := json.NewDecoder(request.Body).Decode(&content); err != nil {
		sendErrorResponse(writer, "Invalid JSON payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields in content (optional step)
	if content.Country == "" || content.IsoCode == "" {
		sendErrorResponse(writer, "Missing required fields: Name or Description", http.StatusBadRequest)
		return
	}

	// Set the current timestamp in RFC3339 format
	timeNow := time.Now().Format(time.RFC3339)
	content.TimeChanged = timeNow

	// Check if Firebase client is initialized
	if firebase.Client == nil {
		sendErrorResponse(writer, "Internal server error: Database client is unavailable.", http.StatusInternalServerError)
		return
	}

	// Write the document to Firestore
	id, _, err := firebase.Client.Collection(collection).Add(firebase.Ctx, content)
	if err != nil {
		sendErrorResponse(writer, "Error when adding document to Firestore: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Document added to collection. Identifier of returned document: " + id.ID)

	// Prepare and send JSON response
	response := consts.RegistrationRequestResponse{
		Id:         id.ID,
		LastChange: timeNow,
	}
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// TODO: Implement the viewDashboardConfiguration function. Issue #6-7
func viewDashboardConfiguration(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNotImplemented)
	writer.Write([]byte("GET method not implemented yet"))
}

// TODO: Implement the replaceDashboardConfiguration function. Issue #8
func replaceDashboardConfiguration(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNotImplemented)
	writer.Write([]byte("PUT method not implemented yet"))
}

// TODO: Implement the deleteDashboardConfiguration function. Issue #9
func deleteDashboardConfiguration(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNotImplemented)
	writer.Write([]byte("DELETE method not implemented yet"))
}

// TODO: !!ADVANCED TASK!! Implement the handleHeadRequest function. Issue #10
func handleHeadRequest(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNotImplemented)
	writer.Write([]byte("HEAD method not implemented yet"))
}

// sendErrorResponse is a helper function to send error responses in JSON format.
func sendErrorResponse(writer http.ResponseWriter, message string, statusCode int) {
	log.Println(message)
	http.Error(writer, message, statusCode)
}
