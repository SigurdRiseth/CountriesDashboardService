package handlers

import (
	"CountriesDashboardService/consts"
	"CountriesDashboardService/utils"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestViewDashboardEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, consts.MockDashboardEndpointWithTestID, nil)
	rec := httptest.NewRecorder()

	// Temporarily override getDashboardConfigFromDB if needed
	GetDashboardConfigFromDBFunc = func(id string) (*consts.RegistrationRequestBody, error) {
		return &consts.RegistrationRequestBody{
			Country: "Norway",
			IsoCode: "NO",
			Features: consts.Features{
				Capital:          utils.BoolPtr(true),
				Coordinates:      utils.BoolPtr(true),
				Population:       utils.BoolPtr(true),
				Area:             utils.BoolPtr(true),
				Temperature:      utils.BoolPtr(true),
				Precipitation:    utils.BoolPtr(true),
				TargetCurrencies: &[]string{"USD", "EUR"},
			},
		}, nil
	}

	ViewDashboard(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	t.Logf("Response: %s", string(body))
}

// TestViewDashboard_InvalidID tests the ViewDashboard function with an invalid ID
func TestViewDashboard_MissingID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, consts.MockDashboardEndpointWithoutID, nil)
	rec := httptest.NewRecorder()

	ViewDashboard(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d", resp.StatusCode)
	}
}

// TestViewDashboard_ValidID tests the ViewDashboard function with an invalid ID
func TestViewDashboard_InvalidID(t *testing.T) {
	GetDashboardConfigFromDBFunc = func(id string) (*consts.RegistrationRequestBody, error) {
		return nil, fmt.Errorf("dashboard not found")
	}

	req := httptest.NewRequest(http.MethodGet, consts.MockDashboardEndpointWithInvalidID, nil)
	rec := httptest.NewRecorder()

	ViewDashboard(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected 500 Internal Server Error, got %d", resp.StatusCode)
	}
}
