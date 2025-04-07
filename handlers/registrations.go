package handlers

import (
	"CountriesDashboardService/consts"
	"CountriesDashboardService/firebase"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
	"net/http"
	"reflect"
	"time"
)

// Collection name in Firestore
const registrationsCollection = "registrations"

// HandleRegistrations handles HTTP requests for dashboard configurations.
// It supports the following methods:
// - POST: Adds a new dashboard configuration.
// - GET: Views all dashboard configurations or a specific one if an {id} parameter is provided.
// - HEAD: Checks if a specific dashboard configuration exists.
// - PUT: Replaces the existing dashboard configuration.
// - DELETE: Deletes the current dashboard configuration.
// - PATCH: Partially updates the existing dashboard configuration.
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
	case http.MethodPatch:
		handlePatchRequest(writer, request)
	default:
		sendErrorResponse(writer, "Unsupported request method: "+request.Method, http.StatusMethodNotAllowed)
		return
	}
}

// addDashboardConfiguration handles HTTP POST requests to add a new dashboard configuration to Firestore.
//
// This function performs the following steps:
// 1. Sets the "Content-Type" header of the response to "application/json".
// 2. Decodes the incoming JSON payload into a `RegistrationRequestBody` struct.
// 3. Adds the current timestamp to the `TimeChanged` field of the decoded struct.
// 4. Writes the decoded and timestamped struct to the Firestore registrationsCollection defined by the `registrationsCollection` constant.
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
	id, _, err := firebase.Client.Collection(registrationsCollection).Add(firebase.Ctx, content)
	if err != nil {
		sendErrorResponse(writer, "Error when adding document to Firestore: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Document added to registrationsCollection. Identifier of returned document: " + id.ID)

	// Prepare and send JSON response
	response := consts.RegistrationRequestResponse{
		Id:         id.ID,
		LastChange: timeNow,
	}
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
	CheckWebhooks(content.IsoCode, Register, id.ID)
}

// TODO: Implement the viewDashboardConfiguration function. Issue #6-7
func viewDashboardConfiguration(writer http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")
	var response interface{}
	var err error

	if id != "" {
		log.Println("Retrieving single dashboard configuration with ID: " + id)
		response, err = GetDashboardConfigFromDBFunc(id) // return a single dashboard configuration
	} else {
		log.Println("Retrieving all dashboard configurations")
		response, err = getAllDashboardConfigsFromDB() // return all dashboard configurations
	}

	if err != nil {
		sendErrorResponse(writer, "Error retrieving dashboard configurations: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode and send the response
	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		sendErrorResponse(writer, "Failed to encode response", http.StatusInternalServerError)
	}
}

// getDashboardConfigFromDB retrieves a specific dashboard configuration from Firestore based on the provided ID.
//
// Parameters:
//   - id: string
//     The Firestore-generated document ID of the dashboard configuration to retrieve.
//
// Returns:
//   - *consts.RegistrationRequestBody: A pointer to the retrieved dashboard configuration struct.
//   - error: An error object if an error occurred during retrieval, otherwise nil.
func getDashboardConfigFromDB(id string) (*consts.RegistrationRequestBody, error) {
	// Retrieve specific message based on id (Firestore-generated hash)
	res := firebase.Client.Collection(registrationsCollection).Doc(id)

	// Retrieve reference to document
	doc, err := res.Get(firebase.Ctx)
	if err != nil {
		log.Println("Error extracting body of returned document of message " + id)
		return nil, err
	}

	var content consts.RegistrationRequestBody
	if err := doc.DataTo(&content); err != nil {
		log.Println("Error converting document data to struct.")
		return nil, err
	}

	return &content, nil
}

// getAllDashboardConfigsFromDB retrieves all dashboard configurations from the Firestore registrationsCollection.
// It returns a slice of RegistrationRequestBody structs or an error if something goes wrong during the process.
//
// Query Details:
//   - Orders the results by the "TimeChanged" field in ascending order.
//   - Limits the retrieval to a predefined maximum number of documents (e.g., 100) for performance reasons.
//   - Gracefully handles Firestore iteration, stopping when all documents have been processed.
//
// Returns:
//   - []consts.RegistrationRequestBody: A slice containing all successfully retrieved and deserialized dashboard configurations.
//   - error: If an error occurs during Firestore document retrieval or data conversion.
//
// Example Usage:
//
//	configs, err := getAllDashboardConfigsFromDB()
//	if err != nil {
//	    log.Fatalf("Failed to get dashboard configurations: %v", err)
//	}
//	fmt.Printf("Retrieved %d configurations.\n", len(configs))
func getAllDashboardConfigsFromDB() ([]consts.RegistrationRequestBody, error) {
	iter := firebase.Client.Collection(registrationsCollection).Limit(100).OrderBy("TimeChanged", firestore.Asc).Documents(firebase.Ctx)

	var content []consts.RegistrationRequestBody
	for doc, err := iter.Next(); !errors.Is(err, iterator.Done); doc, err = iter.Next() {
		if err != nil {
			log.Printf("Error iterating Firestore documents: %v", err)
			return nil, err
		}

		var item consts.RegistrationRequestBody
		if err := doc.DataTo(&item); err != nil {
			log.Printf("Error converting document %s data to struct: %v", doc.Ref.ID, err)
			return nil, err
		}

		content = append(content, item)
	}

	return content, nil
}

// replaceDashboardConfiguration handles HTTP PUT requests to replace an existing dashboard configuration in Firestore.
//
// This function performs the following steps:
// 1. Extracts the configuration ID from the URL path.
// 2. Decodes the incoming JSON payload into a `RegistrationRequestBody` struct.
// 3. Validates the request body for required fields.
// 4. Checks if the configuration exists in Firestore.
// 5. Updates the configuration with the new data, including a new `lastChange` timestamp.
// 6. Returns a 204 No Content response with an empty body on success.
//
// Parameters:
//   - writer: `http.ResponseWriter`
//     The HTTP response writer used to send data back to the client.
//   - request: `*http.Request`
//     The incoming HTTP request, which must include a valid configuration ID in the URL path
//     and a valid JSON payload in the request body.
//
// Behavior:
//   - If no ID is provided in the URL, returns a `400 Bad Request` status.
//   - If the request body is invalid JSON, returns a `400 Bad Request` status.
//   - If required fields are missing, returns a `400 Bad Request` status.
//   - If the Firebase client is not initialized, returns a `500 Internal Server Error` status.
//   - If the document doesn't exist, returns a `404 Not Found` status.
//   - If the update fails, returns a `500 Internal Server Error` status.
//   - Upon successful update, returns a `204 No Content` status with an empty body.
//
// Example Request:
//
//	PUT /dashboard/v1/registrations/516dba7f015f2a68
//	Content-Type: application/json
//	{
//	   "country": "Norway",
//	   "isoCode": "NO",
//	   "features": {
//	                  "temperature": false,
//	                  "precipitation": true,
//	                  "capital": true,
//	                  "coordinates": true,
//	                  "population": true,
//	                  "area": false,
//	                  "targetCurrencies": ["EUR", "SEK"]
//	               }
//	}
func replaceDashboardConfiguration(writer http.ResponseWriter, request *http.Request) {
	log.Println("Processing PUT request for dashboard configuration")

	// Extract the ID from the URL path
	id := request.URL.Query().Get("id")

	// Validate the ID
	if id == "" {
		sendErrorResponse(writer, "Missing configuration ID in URL path", http.StatusBadRequest)
		return
	}

	// Check if Firebase client is initialized
	if firebase.Client == nil {
		sendErrorResponse(writer, "Internal server error: Database client is unavailable", http.StatusInternalServerError)
		return
	}

	// Decode the request body into a RegistrationRequestBody struct
	var content consts.RegistrationRequestBody
	if err := json.NewDecoder(request.Body).Decode(&content); err != nil {
		sendErrorResponse(writer, "Invalid JSON payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if content.Country == "" || content.IsoCode == "" {
		sendErrorResponse(writer, "Missing required fields: country or isoCode", http.StatusBadRequest)
		return
	}

	// Reference to the specific document
	docRef := firebase.Client.Collection(registrationsCollection).Doc(id)

	// Check if document exists first
	doc, err := docRef.Get(firebase.Ctx)
	if err != nil {
		sendErrorResponse(writer, "Error checking document existence: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !doc.Exists() {
		sendErrorResponse(writer, "Configuration not found", http.StatusNotFound)
		return
	}

	// Set the current timestamp in RFC3339 format
	timeNow := time.Now().Format(time.RFC3339)
	content.TimeChanged = timeNow

	// Update the document in Firestore
	_, err = docRef.Set(firebase.Ctx, content)
	if err != nil {
		sendErrorResponse(writer, "Error updating configuration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return 204 No Content on successful update
	writer.WriteHeader(http.StatusNoContent)
	log.Printf("Successfully updated configuration with ID: %s", id)
	CheckWebhooks(content.IsoCode, Change, id)
}

// deleteDashboardConfiguration handles HTTP DELETE requests to remove a specific dashboard configuration from Firestore.
//
// This function performs the following steps:
// 1. Extracts the configuration ID from the URL path.
// 2. Validates that an ID is provided.
// 3. Verifies the document exists before deletion (to distinguish between not-found and other errors).
// 4. Attempts to delete the document from Firestore using the provided ID.
// 5. Returns appropriate HTTP status codes based on the operation's success or failure.
//
// Parameters:
//   - writer: `http.ResponseWriter`
//     The HTTP response writer used to send data back to the client.
//   - request: `*http.Request`
//     The incoming HTTP request, which must include a valid configuration ID in the URL path.
//
// Behavior:
//   - If no ID is provided in the URL, returns a `400 Bad Request` status.
//   - If the Firebase client is not initialized, returns a `500 Internal Server Error` status.
//   - If the document doesn't exist, returns a `404 Not Found` status.
//   - If deletion fails for other reasons, returns a `500 Internal Server Error` status.
//   - Upon successful deletion, returns a `204 No Content` status with an empty body, as per the specification.
//
// Example Request:
//
//	DELETE /dashboard/v1/registrations/?id=516dba7f015f2a68
func deleteDashboardConfiguration(writer http.ResponseWriter, request *http.Request) {
	// Extract the ID from the URL path
	id := request.URL.Query().Get("id")

	// Validate the ID
	if id == "" {
		sendErrorResponse(writer, "Missing configuration ID in URL path", http.StatusBadRequest)
		return
	}

	// Check if Firebase client is initialized
	if firebase.Client == nil {
		sendErrorResponse(writer, "Internal server error: Database client is unavailable", http.StatusInternalServerError)
		return
	}

	// Reference to the specific document
	docRef := firebase.Client.Collection(registrationsCollection).Doc(id)

	// Check if document exists first (to distinguish between not-found and other errors)
	doc, err := docRef.Get(firebase.Ctx)
	if err != nil {
		// Firestore returns a generic error for not found
		sendErrorResponse(writer, "Error checking document existence: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !doc.Exists() {
		sendErrorResponse(writer, "Configuration not found", http.StatusNotFound)
		return
	}

	// Attempt to delete the document
	if _, err := docRef.Delete(firebase.Ctx); err != nil {
		sendErrorResponse(writer, "Error deleting configuration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return 204 No Content on successful deletion
	writer.WriteHeader(http.StatusNoContent)
	log.Printf("Successfully deleted configuration with ID: %s", id)
	CheckWebhooks(doc.Data()["IsoCode"].(string), Delete, id)
}

// handleHeadRequest handles HTTP HEAD requests to check the existence and size of dashboard configurations.
//
// This function performs the following steps:
//  1. Extracts the configuration ID from the URL query parameters.
//  2. If an ID is provided, retrieves the specific dashboard configuration and marshals it to JSON to calculate its size.
//  3. If no ID is provided, retrieves all dashboard configurations and marshals the registrationsCollection to JSON to calculate the total size.
//  4. Sets the "Content-Type" and "Content-Length" headers in the response.
//  5. Returns a 200 OK status with no body to indicate the size of the data.
//
// Parameters:
//   - writer: `http.ResponseWriter`
//     The HTTP response writer used to send headers back to the client.
//   - request: `*http.Request`
//     The incoming HTTP request, which may include a configuration ID in the URL query parameters.
//
// Behavior:
//   - If an error occurs while retrieving a specific configuration or all configurations from the database,
//     it returns a `500 Internal Server Error` status and an error message in the body.
//   - If no error occurs, it marshals the configuration(s) to JSON and calculates the byte size of the data.
//   - Sets the "Content-Length" header to the size of the marshaled JSON data and returns a `200 OK` status
//     to inform the client of the content size (with no response body).
func handleHeadRequest(writer http.ResponseWriter, request *http.Request) {
	// Extract the ID from the URL path
	id := request.URL.Query().Get("id")

	var jsonData []byte
	var err error

	if id != "" {
		// Fetch a specific dashboard config by ID
		doc, err2 := GetDashboardConfigFromDBFunc(id)
		if err2 != nil {
			sendErrorResponse(writer, "Error fetching dashboard config: "+err2.Error(), http.StatusInternalServerError)
			return
		}
		jsonData, err = json.Marshal(doc) // Marshal the single document to JSON
	} else {
		// Fetch and marshal all dashboard configs if no ID is provided
		docs, err2 := getAllDashboardConfigsFromDB()
		if err2 != nil {
			sendErrorResponse(writer, "Error fetching all dashboard configs: "+err2.Error(), http.StatusInternalServerError)
			return
		}
		jsonData, err = json.Marshal(docs) // Marshal the slice of documents to JSON
	}

	if err != nil {
		sendErrorResponse(writer, "Error marshaling data to JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set headers and write response
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(jsonData)+1)) // +1 for newline
	writer.WriteHeader(http.StatusOK)
}

// handlePatchRequest handles HTTP PATCH requests to partially update an existing dashboard configuration in Firestore.
// It extracts the configuration ID from the URL path, validates the request body, checks if the configuration exists,
// applies the partial updates to the document, and updates the 'TimeChanged' timestamp. The function returns a 204 No Content
// response on successful update or appropriate error statuses (400, 404, 500) if validation fails or an error occurs during processing.
//
// The function expects the request body to be a JSON object representing the fields to update, with the 'features' field being
// the main object containing the updatable fields.
//
// Parameters:
//   - writer: http.ResponseWriter - The response writer to send the HTTP response.
//   - request: *http.Request - The incoming HTTP request containing the configuration ID and the update payload.
//
// Behavior:
//   - The configuration ID is extracted from the URL path, and if missing or invalid, a 400 Bad Request is returned.
//   - If the Firestore client is not initialized, a 500 Internal Server Error is returned.
//   - The function checks if the document exists in Firestore, returning a 404 Not Found if the document doesn't exist.
//   - If the request body contains invalid JSON, a 400 Bad Request response is returned.
//   - The 'features' in the payload are iterated, and only non-zero fields are updated in the Firestore document.
//   - The 'TimeChanged' timestamp is updated to the current time.
//   - A 204 No Content response is returned on successful update, with no body in the response.
//
// Example Request:
//
//	PATCH /dashboard/v1/registrations/?id=516dba7f015f2a68
//	Content-Type: application/json
//	{
//	    "features": {
//	        "temperature": false,
//	        "targetCurrencies": ["EUR", "SEK"]
//	    }
//	}
//
// Example Response:
//
//	Success (204 No Content):
//	  HTTP/1.1 204 No Content
//
//	Error (400 Bad Request):
//	  HTTP/1.1 400 Bad Request
//	  {
//	      "error": "Invalid JSON payload: <error details>"
//	  }
//
//	Error (404 Not Found):
//	  HTTP/1.1 404 Not Found
//	  {
//	      "error": "Configuration not found"
//	  }
func handlePatchRequest(writer http.ResponseWriter, request *http.Request) {
	// Extract the ID from the URL path
	id := request.URL.Query().Get("id")

	// Validate the ID
	if id == "" {
		sendErrorResponse(writer, "Missing configuration ID in URL path", http.StatusBadRequest)
		return
	}

	// Check if Firebase client is initialized
	if firebase.Client == nil {
		sendErrorResponse(writer, "Internal server error: Database client is unavailable", http.StatusInternalServerError)
		return
	}

	// Reference to the specific document
	docRef := firebase.Client.Collection(registrationsCollection).Doc(id)
	log.Printf("Patching configuration with ID: %s", id)

	// Check if the document exists
	doc, err := docRef.Get(firebase.Ctx)
	if err != nil {
		sendErrorResponse(writer, "Error checking document existence: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !doc.Exists() {
		sendErrorResponse(writer, "Configuration not found", http.StatusNotFound)
		return
	}

	// Decode the request body into the UserUpdateRequest struct for partial updates
	var inputJSON consts.UserUpdateRequest
	if err := json.NewDecoder(request.Body).Decode(&inputJSON); err != nil {
		sendErrorResponse(writer, "Invalid JSON payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Prepare the updates based on the struct fields using reflection
	var updates []firestore.Update
	val := reflect.ValueOf(inputJSON.Features)
	typ := reflect.TypeOf(inputJSON.Features)

	// Iterate through the fields of the Features struct
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i).Interface()
		// Only add to the update if the value is not nil
		if value != nil {
			updates = append(updates, firestore.Update{Path: "Features." + field.Name, Value: value})
		}
	}

	// Add the lastChange timestamp
	updates = append(updates, firestore.Update{Path: "TimeChanged", Value: time.Now().Format(time.RFC3339)})

	// Perform the update
	_, err = docRef.Update(firebase.Ctx, updates)
	if err != nil {
		sendErrorResponse(writer, "Error updating configuration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return 204 No Content on successful update
	writer.WriteHeader(http.StatusNoContent)
	log.Printf("Successfully patched configuration with ID: %s", id)
	CheckWebhooks(doc.Data()["IsoCode"].(string), Change, id)
}

// sendErrorResponse is a helper function to send error responses in JSON format.
// It logs the error message and sends a response with the specified status code.
func sendErrorResponse(writer http.ResponseWriter, message string, statusCode int) {
	log.Println(message)
	http.Error(writer, message, statusCode)
}
