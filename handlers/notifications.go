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

func deleteWebhook(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Service not implemented"))
}

func retrieveWebhook(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Service not implemented"))
}
