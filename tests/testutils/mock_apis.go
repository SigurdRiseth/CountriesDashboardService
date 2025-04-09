package testutils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

// MockCountriesAPI is a mock server for testing purposes for countries API
func MockCountriesAPI() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{
				"name": map[string]interface{}{
					"common": "Norway",
				},
				"area":       385207,
				"population": 5421241,
				"capital":    []string{"Oslo"},
				"latlng":     []float64{62.0, 10.0},
				"currencies": map[string]interface{}{
					"NOK": map[string]interface{}{
						"name": "Norwegian krone",
					},
				},
			},
		})
	}))
}

// MockMeteoAPI is a mock server for testing weather data
func MockMeteoAPI() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"daily":{"temperature_2m_max":[5],"temperature_2m_min":[-2],"precipitation_sum":[0.15]}}`))
	}))
}

// MockCurrencyAPI is a mock server for testing currency data
func MockCurrencyAPI() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"rates":{"USD":1.2,"EUR":0.9}}`))
	}))
}
