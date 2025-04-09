package handlers

import (
	"CountriesDashboardService/consts"
	"CountriesDashboardService/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func fakeAddDashboardConfiguration(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("addDashboardConfiguration called"))
}

func fakeViewDashboardConfiguration(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("viewDashboardConfiguration called"))
}

func fakeHandleHeadRequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("handleHeadRequest called"))
}

func fakeReplaceDashboardConfiguration(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("replaceDashboardConfiguration called"))
}

func fakeDelete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("deleteDashboardConfiguration called"))
}

func fakeHandlePatchRequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("handlePatchRequest called"))
}

func fakeSendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	http.Error(w, message, statusCode)
}

// TestHandleRegistrations ensures that the central request router delegates to the correct handler
// based on HTTP method (POST, GET, HEAD, PUT, DELETE, PATCH), and handles unsupported methods properly.
func TestHandleRegistrations(t *testing.T) {
	// Inject mocks
	addDashboardConfigurationFunc = fakeAddDashboardConfiguration
	viewDashboardConfigurationFunc = fakeViewDashboardConfiguration
	handleHeadRequestFunc = fakeHandleHeadRequest
	replaceDashboardConfigurationFunc = fakeReplaceDashboardConfiguration
	deleteDashboardConfigurationFunc = fakeDelete
	handlePatchRequestFunc = fakeHandlePatchRequest
	sendErrorResponseFunc = fakeSendErrorResponse

	tests := []struct {
		name         string
		method       string
		expectedCode int
		expectedBody string
	}{
		{"POST method", http.MethodPost, http.StatusOK, "addDashboardConfiguration called"},
		{"GET method", http.MethodGet, http.StatusOK, "viewDashboardConfiguration called"},
		{"HEAD method", http.MethodHead, http.StatusOK, "handleHeadRequest called"},
		{"PUT method", http.MethodPut, http.StatusOK, "replaceDashboardConfiguration called"},
		{"DELETE method", http.MethodDelete, http.StatusOK, "deleteDashboardConfiguration called"},
		{"PATCH method", http.MethodPatch, http.StatusOK, "handlePatchRequest called"},
		{"Unsupported method", http.MethodOptions, http.StatusMethodNotAllowed, "Unsupported request method: OPTIONS"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/registrations", nil)
			rr := httptest.NewRecorder()

			HandleRegistrations(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, rr.Code)
			}

			if strings.TrimSpace(rr.Body.String()) != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, rr.Body.String())
			}
		})
	}
}

// Test_addDashboardConfiguration verifies that a POST request with a valid JSON payload
// results in a successful configuration addition, returning a 201 status and a response with ID and timestamp.
func Test_addDashboardConfiguration(t *testing.T) {
	reqBody := `{
		"country": "Norway",
		"isoCode": "NO",
		"features": {
			"capital": true,
			"coordinates": true,
			"population": true,
			"area": true,
			"temperature": true,
			"precipitation": true,
			"targetCurrencies": ["USD", "EUR"]
		}
	}`

	req := httptest.NewRequest(http.MethodPost, consts.BaseRegistrationPath, strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	//Override Firestore call for isolation.
	AddDashboardConfigToDBFunc = func(body consts.RegistrationRequestBody) (string, error) {
		return "mock-id-123", nil
	}

	addDashboardConfiguration(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201 Created, got %d", resp.StatusCode)
	}

	var jsonResponse map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&jsonResponse); err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	if jsonResponse["id"] == "" || jsonResponse["lastChange"] == "" {
		t.Errorf("Expected 'id' and 'lastChange' in response, got: %+v", jsonResponse)
	}
}

// Test_deleteDashboardConfiguration confirms that a DELETE request with a valid ID
// correctly deletes the config and returns 204 No Content.
func Test_deleteDashboardConfiguration(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/dashboard/v1/registrations/?id=test-id", nil)
	rec := httptest.NewRecorder()

	// Override the deletion logic for isolation
	DeleteDashboardConfigFromDBFunc = func(id string) error {
		if id == "test-id" {
			return nil
		}
		return fmt.Errorf("not found")
	}

	deleteDashboardConfiguration(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("Expected status 204 No Content, got %d", rec.Code)
	}
}

// Test_getAllDashboardConfigsFromDB checks that fetching all dashboard configurations
// from the mock returns the expected dataset without error.
func Test_getAllDashboardConfigsFromDB(t *testing.T) {
	// Setup: mock override
	GetAllDashboardConfigsFunc = func() ([]consts.RegistrationRequestBody, error) {
		return []consts.RegistrationRequestBody{
			{
				Country: "Testland",
				IsoCode: "TL",
				Features: consts.Features{
					Capital:          utils.BoolPtr(true),
					Coordinates:      utils.BoolPtr(true),
					Population:       utils.BoolPtr(true),
					Area:             utils.BoolPtr(true),
					Temperature:      utils.BoolPtr(true),
					Precipitation:    utils.BoolPtr(true),
					TargetCurrencies: &[]string{"USD", "EUR"},
				},
				TimeChanged: "2025-04-07T12:45:00Z",
			},
		}, nil
	}

	got, err := GetAllDashboardConfigsFunc()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("Expected 1 config, got %d", len(got))
	}
	if got[0].Country != "Testland" {
		t.Errorf("Expected country Testland, got %s", got[0].Country)
	}
}

// Test_getDashboardConfigFromDB ensures that retrieving a specific dashboard configuration by ID
// returns the expected structure when the ID is valid.
func Test_getDashboardConfigFromDB(t *testing.T) {
	expected := &consts.RegistrationRequestBody{
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
	}

	GetDashboardConfigFromDBFunc = func(id string) (*consts.RegistrationRequestBody, error) {
		if id != "test-id" {
			t.Errorf("Unexpected ID: got %s, want %s", id, "test-id")
		}
		return expected, nil
	}

	result, err := GetDashboardConfigFromDBFunc("test-id")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

// Test_handleHeadRequest confirms that a HEAD request returns a 200 OK status
// with an empty body, simulating successful endpoint availability checking.
func Test_handleHeadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodHead, consts.BaseRegistrationPath, nil)
	rec := httptest.NewRecorder()

	// Mock the function to return an empty list
	GetAllDashboardConfigsFunc = func() ([]consts.RegistrationRequestBody, error) {
		return []consts.RegistrationRequestBody{}, nil
	}

	handleHeadRequest(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", resp.StatusCode)
	}

	// HEAD responses
	if body := rec.Body.String(); body != "" {
		t.Errorf("Expected empty body for HEAD request, got: %q", body)
	}
}

// Test_handlePatchRequest validates that a PATCH request with valid feature updates and an ID
// leads to a successful partial update and returns a 204 No Content response.
func Test_handlePatchRequest(t *testing.T) {
	reqBody := `{
		"features": {
			"capital": true,
			"coordinates": true
		}
	}`

	// Simulate request with an ID in the query
	req := httptest.NewRequest(http.MethodPatch, consts.BaseRegistrationPath+"?id=test-id", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Mock the patch function
	PatchDashboardConfigInDBFunc = func(id string, input consts.UserUpdateRequest) error {
		if id != "test-id" {
			t.Errorf("Expected ID 'test-id', got '%s'", id)
		}
		if input.Features.Capital == nil || *input.Features.Capital != true {
			t.Errorf("Expected Capital = true, got %+v", input.Features.Capital)
		}
		if input.Features.Coordinates == nil || *input.Features.Coordinates != true {
			t.Errorf("Expected Coordinates = true, got %+v", input.Features.Coordinates)
		}
		return nil
	}

	handlePatchRequest(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204 No Content, got %d", resp.StatusCode)
	}
}

// Test_replaceDashboardConfiguration tests that a full PUT update with valid input and ID
// triggers a replacement in the database and results in a 204 No Content response.
func Test_replaceDashboardConfiguration(t *testing.T) {
	// Request body with valid config
	reqBody := `{
		"country": "Sweden",
		"isoCode": "SE",
		"features": {
			"capital": true,
			"coordinates": true,
			"population": true,
			"area": true,
			"temperature": true,
			"precipitation": true,
			"targetCurrencies": ["USD", "EUR"]
		}
	}`

	// Create test HTTP request
	req := httptest.NewRequest(http.MethodPut, consts.BaseRegistrationPath+"?id=test-id", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Stub to intercept call and verify it gets correct input
	ReplaceDashboardConfigInDBFunc = func(id string, config consts.RegistrationRequestBody) error {
		if id != "test-id" {
			t.Errorf("Expected ID 'test-id', got '%s'", id)
		}
		if config.Country != "Sweden" {
			t.Errorf("Expected country 'Sweden', got '%s'", config.Country)
		}
		if config.IsoCode != "SE" {
			t.Errorf("Expected isoCode 'SE', got '%s'", config.IsoCode)
		}
		if config.Features.TargetCurrencies == nil || len(*config.Features.TargetCurrencies) != 2 {
			t.Errorf("Expected 2 target currencies, got %+v", config.Features.TargetCurrencies)
		}
		return nil // simulate success
	}

	// Call the handler
	replaceDashboardConfiguration(rec, req)

	// Assert result
	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204 No Content, got %d", resp.StatusCode)
	}

}

// Test_replaceDashboardConfiguration_DBError simulates a database failure scenario during a PUT operation,
// verifying that the handler returns a 500 Internal Server Error.
func Test_replaceDashboardConfiguration_DBError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, consts.BaseRegistrationPath+"?id=faulty-id", strings.NewReader(
		`{"country": "X", "isoCode": "X", "features": {}}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ReplaceDashboardConfigInDBFunc = func(id string, config consts.RegistrationRequestBody) error {
		return fmt.Errorf("mock db failure")
	}

	replaceDashboardConfiguration(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected 500 Internal Server Error, got %d", resp.StatusCode)
	}
}

// Test_sendErrorResponse verifies the structure and content of error responses, ensuring
// proper status code, message body, and content type for known errors.
func Test_sendErrorResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	expectedMessage := "Something went wrong"
	expectedStats := http.StatusBadRequest

	sendErrorResponse(rr, expectedMessage, expectedStats)

	if body := strings.TrimSpace(rr.Body.String()); body != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, body)
	}

	if rr.Code != expectedStats {
		t.Errorf("Expected status code %d, got %d", expectedStats, rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/plain; charset=utf-8" {
		t.Errorf("Expected Content-Type 'text/plain; charset=utf-8', got '%s'", contentType)
	}
}

// Test_sendErrorResponse_EmptyMessage confirms that an error response with an empty message
// and 404 status code still sends a compliant response with the correct code.
func Test_sendErrorResponse_EmptyMessage(t *testing.T) {
	rr := httptest.NewRecorder()
	sendErrorResponse(rr, "", http.StatusNotFound)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", rr.Code)
	}
	if strings.TrimSpace(rr.Body.String()) != "" {
		t.Errorf("Expected empty body, got '%s'", rr.Body.String())
	}
}

// Test_viewDashboardConfiguration contains two subtests:
//   - "with ID": Ensures fetching a single config by ID returns the correct result.
//   - "without ID": Verifies that all configs are returned when no ID is provided.
func Test_viewDashboardConfiguration(t *testing.T) {
	t.Run("with ID - fetch single config", func(t *testing.T) {
		expected := &consts.RegistrationRequestBody{
			Country: "Norway",
			IsoCode: "NO",
			Features: consts.Features{
				Capital:     utils.BoolPtr(true),
				Coordinates: utils.BoolPtr(true),
			},
		}

		// Stub single fetch
		GetDashboardConfigFromDBFunc = func(id string) (*consts.RegistrationRequestBody, error) {
			if id != "123" {
				t.Errorf("Expected ID to be '123', got %s", id)
			}
			return expected, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/dashboard/v1/registrations/?id=123", nil)
		rec := httptest.NewRecorder()

		viewDashboardConfiguration(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %d", resp.StatusCode)
		}

		var got consts.RegistrationRequestBody
		if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}

		if !reflect.DeepEqual(&got, expected) {
			t.Errorf("Expected %+v, got %+v", expected, got)
		}
	})

	t.Run("without ID - fetch all configs", func(t *testing.T) {
		expected := []consts.RegistrationRequestBody{
			{
				Country: "Sweden",
				IsoCode: "SE",
				Features: consts.Features{
					Capital:     utils.BoolPtr(true),
					Coordinates: utils.BoolPtr(false),
				},
			},
		}

		GetAllDashboardConfigsFunc = func() ([]consts.RegistrationRequestBody, error) {
			return expected, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/dashboard/v1/registrations/", nil)
		rec := httptest.NewRecorder()

		viewDashboardConfiguration(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
		}

		var got []consts.RegistrationRequestBody
		if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Expected %+v, got %+v", expected, got)
		}
	})
}
