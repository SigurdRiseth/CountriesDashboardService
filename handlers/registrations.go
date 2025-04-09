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
const registrationsCollection = consts.RegistrationsCollection

var AddDashboardConfigToDBFunc = addDashboardConfigurationToFirestore // Override for testing
var DeleteDashboardConfigFromDBFunc = deleteDashboardConfigFromDB
var GetAllDashboardConfigsFunc = getAllDashboardConfigsFromDB
var PatchDashboardConfigInDBFunc = patchDashboardConfigInDB
var ReplaceDashboardConfigInDBFunc = replaceDashboardConfigInDB
var isFirebaseClientInitialized = firebase.IsFirebaseClientInitialized

var (
	addDashboardConfigurationFunc     = addDashboardConfiguration
	viewDashboardConfigurationFunc    = viewDashboardConfiguration
	handleHeadRequestFunc             = handleHeadRequest
	replaceDashboardConfigurationFunc = replaceDashboardConfiguration
	deleteDashboardConfigurationFunc  = deleteDashboardConfiguration
	handlePatchRequestFunc            = handlePatchRequest
	sendErrorResponseFunc             = sendErrorResponse
)

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
		addDashboardConfigurationFunc(writer, request)
	case http.MethodGet:
		viewDashboardConfigurationFunc(writer, request)
	case http.MethodHead:
		handleHeadRequestFunc(writer, request)
	case http.MethodPut:
		replaceDashboardConfigurationFunc(writer, request)
	case http.MethodDelete:
		deleteDashboardConfigurationFunc(writer, request)
	case http.MethodPatch:
		handlePatchRequestFunc(writer, request)
	default:
		sendErrorResponseFunc(writer, "Unsupported request method: "+request.Method, http.StatusMethodNotAllowed)
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
	writer.Header().Set(consts.ContentTypeHeader, consts.ApplicationJSON)

	// Decode the request body into a RegistrationRequestBody struct
	var content consts.RegistrationRequestBody
	if err := json.NewDecoder(request.Body).Decode(&content); err != nil {
		sendErrorResponse(writer, consts.InvalidJSONPayload+err.Error(), http.StatusBadRequest)
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
	//if firebase.client == nil {
	//	sendErrorResponse(writer, "Internal server error: Database client is unavailable.", http.StatusInternalServerError)
	//	return
	//}

	// Write the document to Firestore
	id, err := AddDashboardConfigToDBFunc(content)
	if err != nil {
		sendErrorResponse(writer, "Error when adding document to Firestore: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Document added to registrationsCollection. Identifier of returned document: " + id)
	go CheckWebhooks(content.IsoCode, Register, id)

	// Prepare and send JSON response
	response := consts.RegistrationRequestResponse{
		Id:         id,
		LastChange: timeNow,
	}
	writer.WriteHeader(http.StatusCreated) // 201 Created
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// addDashboardConfigurationToFirestore adds a new dashboard configuration to Firestore.
func addDashboardConfigurationToFirestore(body consts.RegistrationRequestBody) (string, error) {
	if !firebase.IsFirebaseClientInitialized() {
		return "", fmt.Errorf(consts.FBNotInitialized)
	}
	docRef, _, err := firebase.AddToCollection(registrationsCollection, body)
	if err != nil {
		return "", err
	}
	return docRef.ID, nil
}

// viewDashboardConfiguration handles HTTP GET requests to view dashboard configurations.
func viewDashboardConfiguration(writer http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get(consts.QueryParamID)
	var response interface{}
	var err error

	if id != "" {
		log.Println(consts.LogGetSingle + id)
		response, err = GetDashboardConfigFromDBFunc(id) // return a single dashboard configuration
	} else {
		log.Println(consts.LogGetAll)
		response, err = GetAllDashboardConfigsFunc() // return all dashboard configurations
	}

	if err != nil {
		sendErrorResponse(writer, "Error retrieving dashboard configurations: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode and send the response
	writer.Header().Set(consts.ContentTypeHeader, consts.ApplicationJSON)
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
	doc, err := firebase.GetDocument(registrationsCollection, id)
	if err != nil {
		log.Println("Error retrieving document from Firestore: " + err.Error())
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
	iter := firebase.GetLimitedSortedDocuments(registrationsCollection, "TimeChanged", 100)

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
	log.Println(consts.LogPutProcessing)

	id := request.URL.Query().Get(consts.QueryParamID)
	if id == "" {
		sendErrorResponse(writer, consts.MissingIDParamInURL, http.StatusBadRequest)
		return
	}

	var content consts.RegistrationRequestBody
	if err := json.NewDecoder(request.Body).Decode(&content); err != nil {
		sendErrorResponse(writer, consts.InvalidJSONPayload+err.Error(), http.StatusBadRequest)
		return
	}

	if content.Country == "" || content.IsoCode == "" {
		sendErrorResponse(writer, "Missing required fields: country or isoCode", http.StatusBadRequest)
		return
	}

	// Add timestamp before passing to DB layer
	content.TimeChanged = time.Now().Format(time.RFC3339)

	// ReplaceDashboardConfigInDBFunc handles the Firestore logic
	err := ReplaceDashboardConfigInDBFunc(id, content)
	if err != nil {
		sendErrorResponse(writer, "Error updating configuration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	go CheckWebhooks(content.IsoCode, Change, id)

	writer.WriteHeader(http.StatusNoContent)
	log.Printf(consts.LogUpdateSuccess, id)
}

func replaceDashboardConfigInDB(id string, content consts.RegistrationRequestBody) error {
	if !firebase.IsFirebaseClientInitialized() {
		return fmt.Errorf(consts.FBNotInitialized)
	}
	err := firebase.SetDocument(registrationsCollection, id, content)
	return err
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
	id := request.URL.Query().Get(consts.QueryParamID)

	// Validate the ID
	if id == "" {
		sendErrorResponse(writer, consts.MissingIDParamInURL, http.StatusBadRequest)
		return
	}

	// Check if Firebase client is initialized
	if !isFirebaseClientInitialized() {
		sendErrorResponse(writer, "Internal server error: Database client is unavailable", http.StatusInternalServerError)
		return
	}

	isoCode := getIsoCodeFromDocID(id)
	err := DeleteDashboardConfigFromDBFunc(id)
	if err != nil {
		if err.Error() == consts.NotFound {
			sendErrorResponse(writer, "Configuration not found", http.StatusNotFound)
		} else {
			sendErrorResponse(writer, "Error deleting configuration: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	go CheckWebhooks(isoCode, Delete, id)

	// Return 204 No Content on successful deletion
	writer.WriteHeader(http.StatusNoContent)
	log.Printf(consts.LogDeleteSuccess, id)
}

// deleteDashboardConfigFromDB deletes a dashboard configuration from Firestore based on the provided ID.
func deleteDashboardConfigFromDB(id string) error {
	if !firebase.IsFirebaseClientInitialized() {
		return fmt.Errorf(consts.FBNotInitialized)
	}
	docRef := firebase.GetDocumentRef(registrationsCollection, id)

	doc, err := firebase.GetDocumentByRef(docRef)
	if err != nil {
		return err
	}
	if !doc.Exists() {
		return fmt.Errorf(consts.NotFound)
	}

	_, err = firebase.DeleteDocument(docRef)
	return err
}

// handleHeadRequest handles HTTP HEAD requests to check the existence and size of dashboard configurations.
//
// This function performs the following steps:
//  1. Extracts the configuration ID from the URL query parameters.
//  2. If an ID is provided, retrieves the specific dashboard configuration and marshals it to JSON to calculate its size.
//  3. If no ID is provided, retrieves all dashboard configurations and marshals the registrationsCollection to JSON to calculate the total size.
//  4. Sets the "Content-Type" and "Content-Length" headers in the response.
//  5. Returns a 200 OK status without a body to indicate the size of the data.
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
	id := request.URL.Query().Get(consts.QueryParamID)

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
		docs, err2 := GetAllDashboardConfigsFunc()
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
	writer.Header().Set(consts.ContentTypeHeader, consts.ApplicationJSON)
	writer.Header().Set(consts.ContentLengthHeader, fmt.Sprintf("%d", len(jsonData)+1)) // +1 for newline
	writer.WriteHeader(http.StatusOK)
}

// sendErrorResponse is a helper function to send error responses in JSON format.
// It logs the error message and sends a response with the specified status code.
func sendErrorResponse(writer http.ResponseWriter, message string, statusCode int) {
	log.Println(message)
	http.Error(writer, message, statusCode)
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
//   - If the request body contains invalid JSON, a 400 Bad Request response returned.
//   - The 'features' in the payload are iterated, and only non-zero fields updated in the Firestore document.
//   - The 'TimeChanged' timestamp is updated to the current time.
//   - A 204 No Content response is returned on successful update, without body in the response.
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
	id := request.URL.Query().Get(consts.QueryParamID)
	if id == "" {
		sendErrorResponse(writer, consts.MissingIDParamInURL, http.StatusBadRequest)
		return
	}

	var inputJSON consts.UserUpdateRequest
	if err := json.NewDecoder(request.Body).Decode(&inputJSON); err != nil {
		sendErrorResponse(writer, consts.InvalidJSONPayload+err.Error(), http.StatusBadRequest)
		return
	}

	isoCode := getIsoCodeFromDocID(id)

	err := PatchDashboardConfigInDBFunc(id, inputJSON)
	if err != nil {
		if err.Error() == consts.NotFound {
			sendErrorResponse(writer, "Configuration not found", http.StatusNotFound)
		} else {
			sendErrorResponse(writer, "Error updating configuration: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	go CheckWebhooks(isoCode, Change, id)

	writer.WriteHeader(http.StatusNoContent)
	log.Printf(consts.LogPatchSuccess, id)
}

var getIsoCodeFromDocID = func(id string) string {
	doc, err := firebase.GetDocument(registrationsCollection, id)
	if err != nil {
		log.Println("Error retrieving document from Firestore: " + err.Error())
		return ""
	}

	var content consts.RegistrationRequestBody
	if err := doc.DataTo(&content); err != nil {
		log.Println("Error converting document data to struct.")
		return ""
	}

	return content.IsoCode
}

// patchDashboardConfigInDB updates a dashboard configuration in Firestore based on the provided ID and input JSON.
func patchDashboardConfigInDB(id string, inputJSON consts.UserUpdateRequest) error {
	if !firebase.IsFirebaseClientInitialized() {
		return fmt.Errorf(consts.FBNotInitialized)
	}

	docRef := firebase.GetDocumentRef(registrationsCollection, id)

	doc, err := firebase.GetDocumentByRef(docRef)
	if err != nil {
		return fmt.Errorf("error checking document existence: %w", err)
	}
	if !doc.Exists() {
		return fmt.Errorf(consts.NotFound)
	}

	var updates []firestore.Update
	val := reflect.ValueOf(inputJSON.Features)
	typ := reflect.TypeOf(inputJSON.Features)

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i).Interface()
		if value != nil {
			updates = append(updates, firestore.Update{
				Path:  "Features." + field.Name,
				Value: value,
			})
		}
	}

	updates = append(updates, firestore.Update{
		Path:  "TimeChanged",
		Value: time.Now().Format(time.RFC3339),
	})

	err = firebase.UpdateDocument(docRef, updates)
	return err
}
