package testutils

import (
	"net/http"
	"net/http/httptest"
)

// MockCountriesAPI creates a mock server for the countries API
func MockCountriesAPI() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
}

// MockMeteoAPI creates a mock server for the meteo API
func MockMeteoAPI() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
}

// MockCurrencyAPI creates a mock server for the currency API
func MockCurrencyAPI() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
}
