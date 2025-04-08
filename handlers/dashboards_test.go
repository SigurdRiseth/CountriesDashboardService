package handlers

import (
	"CountriesDashboardService/consts"
	"CountriesDashboardService/tests/testutils"
	"CountriesDashboardService/utils"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	CheckWebhooks = testutils.MockCheckWebhooks
}

func TestViewDashboard_Success(t *testing.T) {
	// Stub DB call
	GetDashboardConfigFromDBFunc = func(id string) (*consts.RegistrationRequestBody, error) {
		return &consts.RegistrationRequestBody{
			Country: "Norway",
			IsoCode: "NO",
			Features: consts.Features{
				Temperature:      utils.BoolPtr(true),
				Precipitation:    utils.BoolPtr(true),
				Capital:          utils.BoolPtr(true),
				Coordinates:      utils.BoolPtr(true),
				Population:       utils.BoolPtr(true),
				Area:             utils.BoolPtr(true),
				TargetCurrencies: &[]string{"USD", "EUR"},
			},
		}, nil
	}

	// Stub external APIs
	FetchCountryDataFunc = func(country string) map[string]interface{} {
		return map[string]interface{}{
			"capital":    []interface{}{"Oslo"},
			"latlng":     []interface{}{59.91, 10.75},
			"population": float64(5379475),
			"area":       float64(323802),
			"currencies": map[string]interface{}{"NOK": map[string]interface{}{}},
		}
	}

	FetchWeatherDataFunc = func(features map[string]interface{}, lat, lon float64, config *consts.RegistrationRequestBody) {
		features["temperature"] = -1.2
		features["precipitation"] = 0.8
	}

	FetchCurrencyDataFunc = func(features map[string]interface{}, countryData map[string]interface{}, targetCurrencies []string) {
		features["targetCurrencies"] = map[string]interface{}{
			"USD": 1.05,
			"EUR": 0.95,
		}
	}

	// Prepare request
	req := httptest.NewRequest(http.MethodGet, "/dashboard/v1/dashboards?id=test-id", nil)
	rr := httptest.NewRecorder()

	ViewDashboard(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("response is not valid JSON: %v", err)
	}

	// Check expected fields
	if response["country"] != "Norway" {
		t.Errorf("unexpected country: got %v", response["country"])
	}
	if _, ok := response["features"].(map[string]interface{})["temperature"]; !ok {
		t.Errorf("missing temperature in features")
	}
}

func TestViewDashboard_MissingID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/dashboard/v1/dashboards", nil)
	rr := httptest.NewRecorder()

	ViewDashboard(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", rr.Code)
	}
}

func TestViewDashboard_DBError(t *testing.T) {
	GetDashboardConfigFromDBFunc = func(id string) (*consts.RegistrationRequestBody, error) {
		return nil, errors.New("DB error")
	}

	req := httptest.NewRequest(http.MethodGet, "/dashboard/v1/dashboards?id=test-id", nil)
	rr := httptest.NewRecorder()

	ViewDashboard(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 Internal Server Error, got %d", rr.Code)
	}
}

func TestViewDashboard_FailedCountryData(t *testing.T) {
	GetDashboardConfigFromDBFunc = func(id string) (*consts.RegistrationRequestBody, error) {
		return &consts.RegistrationRequestBody{
			Country: "Norwy",
			IsoCode: "N",
			Features: consts.Features{
				Capital: utils.BoolPtr(true),
			},
		}, nil
	}

	FetchCountryDataFunc = func(_ string) map[string]interface{} {
		return nil
	}

	req := httptest.NewRequest(http.MethodGet, "/dashboard/v1/dashboards?id=test-id", nil)
	rr := httptest.NewRecorder()

	ViewDashboard(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 Internal Server Error, got %d", rr.Code)
	}
}
