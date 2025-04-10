package cache

import (
	"CountriesDashboardService/consts"
	"fmt"
)

// Centralized cache key builders

// CountryCacheKey generates a cache key for country information
func CountryCacheKey(isoOrName string) string {
	return fmt.Sprintf(consts.CountryKeyStamp, isoOrName)
}

// WeatherCacheKey generates a cache key for weather information
func WeatherCacheKey(lat, lon float64) string {
	return fmt.Sprintf(consts.WeatherKeyStamp, lat, lon)
}

// CurrencyCacheKey generates a cache key for currency rates
func CurrencyCacheKey(baseCurrency string) string {
	return fmt.Sprintf(consts.CurrencyKeyStamp, baseCurrency)
}
