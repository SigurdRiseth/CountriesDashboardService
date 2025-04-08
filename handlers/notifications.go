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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net/http"
	"time"
)

// Production function call
var addToCollection = firebase.AddToCollection
var getDocument = firebase.GetDocument
var getDocumentByRef = firebase.GetDocumentByRef
var getDocumentRef = firebase.GetDocumentRef
var deleteDocument = firebase.DeleteDocument
var getCollectionIterator = firebase.GetCollectionIterator
var firebaseClientInitialized = firebase.FirebaseClientInitialized
var exists = firebase.DocumentExists

// Constants (assumed to be declared elsewhere)
var (
	Secret       = []byte("your-secret-key")
	SignatureKey = "X-Signature"
	maxRetries   = 3
)

const (
	Register = "REGISTER"
	Change   = "CHANGE"
	Delete   = "DELETE"
	Invoke   = "INVOKE"
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

	// Check if body is nil
	if request.Body == nil {
		sendErrorResponse(writer, "Empty request body", http.StatusBadRequest)
		return
	}

	// Decode the request body into a WebhookRegistration struct
	var webhook consts.WebhookRegistration
	if err := json.NewDecoder(request.Body).Decode(&webhook); err != nil {
		sendErrorResponse(writer, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validate the Event field
	validEvents := map[string]bool{Register: true, Change: true, Delete: true, Invoke: true}
	if !validEvents[webhook.Event] {
		sendErrorResponse(writer, "Invalid Event value: "+webhook.Event, http.StatusBadRequest)
		return
	}

	// Set the current timestamp in RFC3339 format
	timeNow := time.Now().Format(time.RFC3339)
	webhook.TimeChanged = timeNow

	// Write the document to Firestore
	id, _, err := addToCollection(notificationsCollection, webhook)
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
	if !firebaseClientInitialized() {
		sendErrorResponse(writer, "Internal server error: Database client is unavailable", http.StatusInternalServerError)
		return
	}

	// Reference to the specific document
	docRef := getDocumentRef(notificationsCollection, id)

	// Check if document exists first (to distinguish between not-found and other errors)
	doc, err := getDocumentByRef(docRef)
	if err != nil {
		// Firestore returns a generic error for not found
		sendErrorResponse(writer, "Error checking document existence: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists(doc) {
		sendErrorResponse(writer, "Configuration not found", http.StatusNotFound)
		return
	}

	// Attempt to delete the document
	if _, err := deleteDocument(docRef); err != nil {
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
	// Retrieve reference to document
	doc, err := getDocument(notificationsCollection, id)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			sendErrorResponse(writer, "Webhook not found", http.StatusNotFound)
			return
		}
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
	iter := getCollectionIterator(notificationsCollection)

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

// CheckWebhooks checks for registered webhooks that match the given country and event,
// and triggers a webhook call for each matching registration.
var CheckWebhooks = func(country, event, id string) {
	timeStamp := time.Now().Format(time.RFC3339)

	// Retrieve all webhook registrations from Firestore
	iter := getCollectionIterator(notificationsCollection)

	var webhooks []consts.WebhookRegistration
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break // Exit loop when all documents are processed
		}
		if err != nil {
			log.Printf("Error iterating Firestore documents: %v", err)
			return
		}

		var webhook consts.WebhookRegistration
		if err := doc.DataTo(&webhook); err != nil {
			log.Printf("Error converting Firestore document to struct: %v", err)
			continue // Skip this document and proceed with the next one
		}

		webhooks = append(webhooks, webhook)
	}

	// Filter and trigger matching webhooks
	for _, webhook := range webhooks {
		if webhook.Event == event && webhook.Country == country {
			log.Printf("Triggering webhook for country %s with event %s: %s", country, event, webhook.Url)
			CallUrl(webhook.Url, event, id, timeStamp, webhook.Country)
		}
	}
}

// CallUrl sends an HTTP POST request with event details, including HMAC-based validation for payload security.
func CallUrl(url, event, id, timeStamp, country string) {
	log.Printf("Attempting invocation of URL %s with event '%s' and ID '%s'.", url, event, id)

	// Create the JSON payload
	payload := consts.WebhookPayload{
		Event:   event,
		Id:      id,
		Time:    timeStamp,
		Country: country,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling payload: %v", err)
		return
	}

	// Create the HTTP POST request with JSON body and content-based HMAC header
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payloadBytes))
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(SignatureKey, generateHMACSignature(payloadBytes))

	// Perform the HTTP POST request with retry logic
	client := &http.Client{Timeout: 5 * time.Second}
	executeWithRetry(client, req, url)
}

// generateHMACSignature generates an HMAC signature for the given payload.
func generateHMACSignature(payload []byte) string {
	mac := hmac.New(sha256.New, Secret)
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

// executeWithRetry performs the HTTP POST request with retry and exponential backoff.
func executeWithRetry(client *http.Client, req *http.Request, url string) {
	for retry := 0; retry < maxRetries; retry++ {
		res, err := client.Do(req)
		if err != nil {
			log.Printf("HTTP request failed (attempt %d/%d): %v", retry+1, maxRetries, err)
			time.Sleep(time.Duration(retry) * time.Second)
			continue
		}

		defer res.Body.Close()

		// Read and log the response
		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			return
		}

		log.Printf("Webhook %s invoked. Received status code %d and body: %s", url, res.StatusCode, string(responseBody))

		// Successful response handling
		if res.StatusCode >= 200 && res.StatusCode < 300 {
			log.Printf("Successfully invoked webhook %s with status code %d", url, res.StatusCode)
			return
		}

		log.Printf("Non-2xx status code received: %d", res.StatusCode)
		if retry == maxRetries-1 {
			log.Printf("Max retries reached. Failed to invoke webhook %s.", url)
		}
		time.Sleep(time.Duration(retry+1) * time.Second) // Exponential backoff before retry
	}

	log.Println("Webhook invocation failed after retries.")
}
