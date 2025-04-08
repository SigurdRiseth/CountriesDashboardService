package handlers

import (
	"CountriesDashboardService/consts"
	"CountriesDashboardService/tests/testutils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

func TestRegisteringValidJSON(t *testing.T) {

	addToCollection = testutils.MockAddToCollection

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

func TestDeleteWebhookInvalidID(t *testing.T) {
	// Mock the deleteDocument function
	deleteDocument = testutils.MockDeleteDocument

	// Create a request with a invalid ID
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

func TestDeleteWebhookValidID(t *testing.T) {
	// Mock the deleteDocument function
	deleteDocument = testutils.MockDeleteDocument
	getDocumentRef = testutils.MockGetDocumentRef
	getDocumentByRef = testutils.MockGetDocumentByRef
	firebaseClientInitialized = testutils.MockFirebaseClientInitialized
	exists = testutils.MockDocumentExists

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
