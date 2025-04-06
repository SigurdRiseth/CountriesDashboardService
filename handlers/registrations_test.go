package handlers

import (
	"CountriesDashboardService/consts"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestHandleRegistrations(t *testing.T) {
	type args struct {
		writer  http.ResponseWriter
		request *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleRegistrations(tt.args.writer, tt.args.request)
		})
	}
}

// Test POST request with valid payload → should return 201 Created with document ID and timestamp.
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

	req := httptest.NewRequest(http.MethodPost, consts.RegistrationEndpoint, strings.NewReader(reqBody))
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

func Test_deleteDashboardConfiguration(t *testing.T) {
	type args struct {
		writer  http.ResponseWriter
		request *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deleteDashboardConfiguration(tt.args.writer, tt.args.request)
		})
	}
}

func Test_getAllDashboardConfigsFromDB(t *testing.T) {
	tests := []struct {
		name    string
		want    []consts.RegistrationRequestBody
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getAllDashboardConfigsFromDB()
			if (err != nil) != tt.wantErr {
				t.Errorf("getAllDashboardConfigsFromDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAllDashboardConfigsFromDB() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getDashboardConfigFromDB(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    *consts.RegistrationRequestBody
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDashboardConfigFromDB(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDashboardConfigFromDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDashboardConfigFromDB() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleHeadRequest(t *testing.T) {
	type args struct {
		writer  http.ResponseWriter
		request *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handleHeadRequest(tt.args.writer, tt.args.request)
		})
	}
}

func Test_handlePatchRequest(t *testing.T) {
	type args struct {
		writer  http.ResponseWriter
		request *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlePatchRequest(tt.args.writer, tt.args.request)
		})
	}
}

func Test_replaceDashboardConfiguration(t *testing.T) {
	type args struct {
		writer  http.ResponseWriter
		request *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			replaceDashboardConfiguration(tt.args.writer, tt.args.request)
		})
	}
}

func Test_sendErrorResponse(t *testing.T) {
	type args struct {
		writer     http.ResponseWriter
		message    string
		statusCode int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendErrorResponse(tt.args.writer, tt.args.message, tt.args.statusCode)
		})
	}
}

func Test_viewDashboardConfiguration(t *testing.T) {
	type args struct {
		writer  http.ResponseWriter
		request *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viewDashboardConfiguration(tt.args.writer, tt.args.request)
		})
	}
}
