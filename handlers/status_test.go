package handlers

import (
	"CountriesDashboardService/consts"
	"CountriesDashboardService/utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// mockCountWebhooks simulates returning 5 registered webhooks
func mockCountWebhooks() (int64, error) {
	return 5, nil
}

// mockNotificationDB simulates Firestore DB is accessible
func mockNotificationDB() int {
	return http.StatusOK
}

func TestStatusWithMockedServices(t *testing.T) {
	// Set a static start time for uptime testing
	start := time.Now().Add(-time.Hour) // simulate 1 hour uptime
	utils.StartTime = start

	// Inject mocks
	CountWebhooksFunc = mockCountWebhooks
	CheckNotificationDBFunc = mockNotificationDB

	// Create test request
	req := httptest.NewRequest(http.MethodGet, consts.BaseStatusPath, nil)
	w := httptest.NewRecorder()

	// Call handler
	HandleStatus(w, req)

	// Check response status
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	// Decode JSON response
	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Required fields to validate
	requiredFields := []string{
		"countries_api", "meteo_api", "currency_api",
		"notification_db", "webhooks", "version", "uptime",
	}

	for _, key := range requiredFields {
		if _, ok := resp[key]; !ok {
			t.Errorf("Missing field in response: %s", key)
		}
	}

	// Additional type checks
	if resp["version"] != "v1" {
		t.Errorf("Expected version to be 'v1', got %v", resp["version"])
	}
	if resp["webhooks"] != float64(5) {
		t.Errorf("Expected 5 webhooks, got %v", resp["webhooks"])
	}
	if resp["notification_db"] != float64(200) {
		t.Errorf("Expected notification_db to be 200, got %v", resp["notification_db"])
	}
}

func TestStatusInvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, consts.BaseStatusPath, nil)
	w := httptest.NewRecorder()

	HandleStatus(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Method not allowed") {
		t.Errorf("Expected error message in response, got %s", w.Body.String())
	}
}
