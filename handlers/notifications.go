package handlers

import (
	"CountriesDashboardService/consts"
	"CountriesDashboardService/firebase"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"google.golang.org/api/iterator"
	"io"
	"log"
	"net/http"
	"time"
)

// Constants (assumed to be declared elsewhere)
var (
	Secret       = []byte("your-secret-key")
	SignatureKey = "X-Signature"
	maxRetries   = 3
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
		sendErrorResponse(writer, "Error encoding JSON response: "+err.Error(), http.StatusInternalServerError)
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

// retrieveWebhook handles the retrieval of webhook configurations.
// It checks if an ID is provided in the URL query parameters and calls
// the appropriate function to either retrieve a specific webhook or all webhooks.
//
// Parameters:
// - writer: http.ResponseWriter to write the response
// - request: *http.Request containing the request data
func retrieveWebhook(writer http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")

	if id != "" {
		retrieveSpecificWebhook(writer, request, id)
	} else {
		retrieveAllWebhooks(writer, request)
	}
}

// retrieveSpecificWebhook handles the retrieval of a specific webhook configuration.
// It retrieves the document from Firestore based on the provided ID, converts it to
// a WebhookRegistration struct, and sends it as a JSON response.
//
// Parameters:
// - writer: http.ResponseWriter to write the response
// - request: *http.Request containing the request data
// - id: string representing the Firestore document ID
func retrieveSpecificWebhook(writer http.ResponseWriter, request *http.Request, id string) {
	// Retrieve specific message based on id (Firestore-generated hash)
	res := firebase.Client.Collection(notificationsCollection).Doc(id)

	// Retrieve reference to document
	doc, err := res.Get(firebase.Ctx)
	if err != nil {
		sendErrorResponse(writer, "Error extracting body of returned document of message "+id, http.StatusInternalServerError)
		return
	}

	var content consts.WebhookRegistration
	if err := doc.DataTo(&content); err != nil {
		sendErrorResponse(writer, "Error converting document data to struct.", http.StatusInternalServerError)
		return
	}

	// Send response
	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(content); err != nil {
		sendErrorResponse(writer, "Error encoding JSON response.", http.StatusInternalServerError)
	}
}

// retrieveAllWebhooks handles the retrieval of all webhook configurations.
// It iterates over all documents in the Firestore collection, converts them
// to WebhookRegistration structs, and sends them as a JSON response.
//
// Parameters:
// - writer: http.ResponseWriter to write the response
// - request: *http.Request containing the request data
func retrieveAllWebhooks(writer http.ResponseWriter, request *http.Request) {
	// Retrieve all messages
	iter := firebase.Client.Collection(notificationsCollection).Documents(firebase.Ctx)

	// Prepare a slice to hold all messages
	var messages []consts.WebhookRegistration

	for doc, err := iter.Next(); !errors.Is(err, iterator.Done); doc, err = iter.Next() {
		if err != nil {
			sendErrorResponse(writer, "Error iterating Firestore documents.", http.StatusInternalServerError)
			return
		}

		var item consts.WebhookRegistration
		if err := doc.DataTo(&item); err != nil {
			sendErrorResponse(writer, "Error converting document data to struct.", http.StatusInternalServerError)
			return
		}

		messages = append(messages, item)
	}

	// Send response
	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(messages); err != nil {
		sendErrorResponse(writer, "Error encoding JSON response.", http.StatusInternalServerError)
	}
}

// CallUrl sends an HTTP POST request with event and content, using HMAC validation.
//
// Parameters:
// - url: The URL to which the HTTP POST request is sent.
// - event: The event type to be included in the payload.
// - content: The content to be included in the payload.
func CallUrl(url, event, content string) {
	log.Printf("Attempting invocation of URL %s with event '%s' and content: '%s'.", url, event, content)

	// Create the JSON payload
	payload := consts.WebhookPayload{
		Event:   event,
		Content: content,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling payload: %v", err)
		return
	}

	// Create the HTTP POST request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payloadBytes))
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// HMAC signature generation
	mac := hmac.New(sha256.New, Secret)
	_, err = mac.Write([]byte(content))
	if err != nil {
		log.Printf("Error during content hashing: %v", err)
		return
	}
	req.Header.Set(SignatureKey, hex.EncodeToString(mac.Sum(nil)))

	// Perform the HTTP POST request with retry logic
	client := &http.Client{Timeout: 5 * time.Second}
	for retry := 0; retry < maxRetries; retry++ {
		res, err := client.Do(req)
		if err != nil {
			log.Printf("HTTP request failed (attempt %d/%d): %v", retry+1, maxRetries, err)
			time.Sleep(time.Duration(retry) * time.Second) // Exponential backoff
			continue
		}

		defer res.Body.Close()

		// Read the response body
		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			return
		}

		log.Printf("Webhook %s invoked. Received status code %d and body: %s", url, res.StatusCode, string(responseBody))

		// Check for HTTP errors and decide whether to retry
		if res.StatusCode >= 200 && res.StatusCode < 300 {
			log.Printf("Successfully invoked webhook %s with status code %d", url, res.StatusCode)
			return
		} else {
			log.Printf("Non-2xx status code received: %d", res.StatusCode)
			if retry == maxRetries-1 {
				log.Printf("Max retries reached. Failed to invoke webhook %s.", url)
				return
			}
			time.Sleep(time.Duration(retry+1) * time.Second) // Exponential backoff before retry
		}
	}

	log.Println("Webhook invocation failed after retries.")
}
