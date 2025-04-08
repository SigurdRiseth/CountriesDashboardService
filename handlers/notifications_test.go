package handlers

import (
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
