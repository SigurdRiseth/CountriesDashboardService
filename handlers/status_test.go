package handlers

import (
	"CountriesDashboardService/consts"
	"CountriesDashboardService/tests/testutils"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestStatusWithMockedServices tests the status handler with mocked external services
func TestStatusWithMockedServices(t *testing.T) {
	// Start mock servers
	countriesMock := testutils.MockCountriesAPI()
	defer countriesMock.Close()

	meteoMock := testutils.MockMeteoAPI()
	defer meteoMock.Close()

	currencyMock := testutils.MockCurrencyAPI()
	defer currencyMock.Close()

	// Override API endpoints
	consts.RestCountriesAPITest = countriesMock.URL
	consts.OpenMeteoAPITest = meteoMock.URL
	consts.CurrencyAPITest = currencyMock.URL

	// Mock the firebase functions
	CountWebhooksFunc = func() (int64, error) {
		return 5, nil // 5 as dummy value
	}

	// Mock notification DB check to return 200 OK
	CheckNotificationDBFunc = func() int {
		return http.StatusOK
	}

	// Make test request
	req := httptest.NewRequest(http.MethodGet, "/dashboard/v1/status/", nil)
	rr := httptest.NewRecorder()

	handleGetStatus(rr, req)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	t.Logf("Response: %s", rr.Body.String())
}
