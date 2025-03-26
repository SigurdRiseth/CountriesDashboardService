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
const notificationsCollection = "notifications"

func HandleNotifications(writer http.ResponseWriter, request *http.Request) {
	log.Println("Notifications endpoint received a " + request.Method + " request")
	switch request.Method {
	case http.MethodPost:
		registerWebhook(writer, request)
	case http.MethodDelete:
		deleteWebhook(writer, request)
	case http.MethodGet:
		retrieveWebhook(writer, request)
	default:
		sendErrorResponse(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// registerWebhook handles the registration of a new webhook.
// It decodes the request body into a WebhookRegistration struct,
// validates the Event field, sets the current timestamp, writes
// the document to Firestore, and sends a JSON response.
//
// Parameters:
// - writer: http.ResponseWriter to write the response
// - request: *http.Request containing the request data
func registerWebhook(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	// Decode the request body into a WebhookRegistration struct
	var webhook consts.WebhookRegistration
	if err := json.NewDecoder(request.Body).Decode(&webhook); err != nil {
		sendErrorResponse(writer, "Invalid JSON payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the Event field
	validEvents := map[string]bool{"REGISTER": true, "CHANGE": true, "DELETE": true, "INVOKE": true}
	if !validEvents[webhook.Event] {
		sendErrorResponse(writer, "Invalid Event value: "+webhook.Event, http.StatusBadRequest)
		return
	}

	// Set the current timestamp in RFC3339 format
	timeNow := time.Now().Format(time.RFC3339)
	webhook.TimeChanged = timeNow

	// Write the document to Firestore
	id, _, err := firebase.Client.Collection(notificationsCollection).Add(firebase.Ctx, webhook)
	if err != nil {
		sendErrorResponse(writer, "Error when adding document to Firestore: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare and send JSON response
	response := consts.RegistrationRequestResponse{
		Id:         id.ID,
		LastChange: timeNow,
	}
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// deleteWebhook handles the deletion of a webhook configuration.
// It extracts the ID from the URL path, validates it, checks if the document exists,
// and attempts to delete it from Firestore. It sends appropriate HTTP responses
// based on the outcome.
//
// Parameters:
// - writer: http.ResponseWriter to write the response
// - request: *http.Request containing the request data
func deleteWebhook(writer http.ResponseWriter, request *http.Request) {
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
	docRef := firebase.Client.Collection(notificationsCollection).Doc(id)

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
}

func retrieveWebhook(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Service not implemented"))
}
