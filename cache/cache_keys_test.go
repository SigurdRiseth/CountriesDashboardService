package cache

import (
	"testing"
)

func TestCountryCacheKey(t *testing.T) {
	got := CountryCacheKey("NO")
	want := "country_NO"
	if got != want {
		t.Errorf("CountryCacheKey() = %q, want %q", got, want)
	}
}

func TestWeatherCacheKey(t *testing.T) {
	got := WeatherCacheKey(59.91, 10.75)
	want := "weather_59.91_10.75"
	if got != want {
		t.Errorf("WeatherCacheKey() = %q, want %q", got, want)
	}
}

func TestCurrencyCacheKey(t *testing.T) {
	got := CurrencyCacheKey("USD")
	want := "currency_USD"
	if got != want {
		t.Errorf("CurrencyCacheKey() = %q, want %q", got, want)
	}
}
