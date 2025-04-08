package handlers

import (
	"CountriesDashboardService/consts"
	"CountriesDashboardService/tests/testutils"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	// Mock the firestore functions to avoid actual Firestore calls
	deleteDocument = testutils.MockDeleteDocument
	getDocumentRef = testutils.MockGetDocumentRef
	getDocumentByRef = testutils.MockGetDocumentByRef
	firebaseClientInitialized = testutils.MockFirebaseClientInitialized
	exists = testutils.MockDocumentExists
	addToCollection = testutils.MockAddToCollection
	getDocument = testutils.MockGetDocument
}

// TestIllegalRequest tests the HandleNotifications function for an unsupported HTTP method.
func TestIllegalRequest(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, "/dashboard/v1/notifications/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the HandleNotifications function with the request and recorder
	HandleNotifications(rr, req)

	// Check the status code
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusMethodNotAllowed)
	}
}

// TestRegisteringValidJSON tests the HandleNotifications function for a valid JSON body.
func TestRegisteringValidJSON(t *testing.T) {

	// Create a request with a valid JSON body
	req, err := http.NewRequest(http.MethodPost, "/dashboard/v1/notifications/", bytes.NewReader([]byte(`{"event":"DELETE"}`)))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the HandleNotifications function with the request and recorder
	HandleNotifications(rr, req)

	// Decode the response body into a RegistrationRequestResponse struct
	var response consts.RegistrationRequestResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Assert that the id field matches the expected value
	if response.Id != "mockID" {
		t.Errorf("Handler returned wrong id: got %v want %v", response.Id, "mockID")
	}

	// Assert that the last change field is not empty
	if response.LastChange == "" {
		t.Errorf("Handler returned empty last change: got %v want non-empty", response.LastChange)
	}

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}
}

// TestRegisteringIllegalJSON tests the HandleNotifications function for an invalid JSON body.
func TestRegisteringIllegalJSON(t *testing.T) {
	// Create a request with an invalid JSON body
	req, err := http.NewRequest(http.MethodPost, "/dashboard/v1/notifications/", bytes.NewBufferString("{invalid_json}"))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the HandleNotifications function with the request and recorder
	HandleNotifications(rr, req)

	// Check the status code
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusBadRequest)
	}
}

// TestRegisteringEmptyJSON tests the HandleNotifications function for an empty JSON body.
func TestRegisteringEmptyJSON(t *testing.T) {
	// Create a request with an empty JSON body
	req, err := http.NewRequest(http.MethodPost, "/dashboard/v1/notifications/", bytes.NewBufferString("{}"))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the HandleNotifications function with the request and recorder
	HandleNotifications(rr, req)

	// Check the status code
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusBadRequest)
	}
}

// TestRegisteringNilJSON tests the HandleNotifications function for a nil JSON body.
func TestRegisteringNilJSON(t *testing.T) {
	// Create a request with a nil JSON body
	req, err := http.NewRequest(http.MethodPost, "/dashboard/v1/notifications/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the HandleNotifications function with the request and recorder
	HandleNotifications(rr, req)

	// Check the status code
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusBadRequest)
	}
}

// TestRegisteringIllegalEvent tests the HandleNotifications function for an invalid event type.
func TestRegisteringIllegalEvent(t *testing.T) {
	// Create a request with an invalid event type
	req, err := http.NewRequest(http.MethodPost, "/dashboard/v1/notifications/", bytes.NewReader([]byte(`{"event":"illegal_event"}`)))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the HandleNotifications function with the request and recorder
	HandleNotifications(rr, req)

	// Check the status code
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusBadRequest)
	}
}

// TestDeleteWebhookInvalidID tests the HandleNotifications function for a DELETE request with an invalid ID.
func TestDeleteWebhookInvalidID(t *testing.T) {
	// Create a request with an invalid ID
	req, err := http.NewRequest(http.MethodDelete, "/dashboard/v1/notifications/?id=", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the HandleNotifications function with the request and recorder
	HandleNotifications(rr, req)

	// Check the status code
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusBadRequest)
	}
}

// TestDeleteWebhookNoID tests the HandleNotifications function for a DELETE request without an ID.
func TestDeleteWebhookNoID(t *testing.T) {
	// Mock the deleteDocument function
	deleteDocument = testutils.MockDeleteDocument

	// Create a request with an invalid ID
	req, err := http.NewRequest(http.MethodDelete, "/dashboard/v1/notifications/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the HandleNotifications function with the request and recorder
	HandleNotifications(rr, req)

	// Check the status code
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusBadRequest)
	}
}

// TestDeleteWebhookValidID tests the HandleNotifications function for a DELETE request with a valid ID.
func TestDeleteWebhookValidID(t *testing.T) {

	// Create a request with a valid ID
	req, err := http.NewRequest(http.MethodDelete, "/dashboard/v1/notifications/?id=mockID", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the HandleNotifications function with the request and recorder
	HandleNotifications(rr, req)

	// Check the status code
	if rr.Code != http.StatusNoContent {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusNoContent)
	}
}

// TestGetWebhookUnusedID tests the HandleNotifications function for a GET request with an unused ID.
func TestGetWebhookUnusedID(t *testing.T) {
	// Create a request with an unused ID
	req, err := http.NewRequest(http.MethodGet, "/dashboard/v1/notifications/?id=unused", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the HandleNotifications function with the request and recorder
	HandleNotifications(rr, req)

	// Check the status code
	if rr.Code != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusNotFound)
	}
}

// TestCallUrl tests the CallUrl function for a successful HTTP POST request.
func TestCallUrl(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Verify the content type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Verify the HMAC signature
		body, _ := io.ReadAll(r.Body)
		expectedSignature := generateHMACSignature(body)
		if r.Header.Get(SignatureKey) != expectedSignature {
			t.Errorf("Invalid HMAC signature: got %s, want %s", r.Header.Get(SignatureKey), expectedSignature)
		}

		// Respond with a success status
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	// Call the function with test data
	CallUrl(mockServer.URL, "DELETE", "testID", "2023-01-01T00:00:00Z", "TestCountry")
}
