package cache

import (
	"CountriesDashboardService/consts"
	"CountriesDashboardService/firebase"
	"context"
	"fmt"
	"time"
)

// CacheItem defines a cached Firestore document
type CacheItem struct {
	Timestamp time.Time   `firestore:"timestamp"`
	Data      interface{} `firestore:"data"`
}

// isCacheExpired returns true if the given timestamp is older than maxAge
func isCacheExpired(ts time.Time, maxAge time.Duration) bool {
	return time.Since(ts) > maxAge
}

// setCache stores data in the cache with a timestamp
func setCache(key string, data interface{}) error {
	if firebase.Client == nil {
		return fmt.Errorf(consts.FTNotInitialized)
	}
	_, err := firebase.Client.Collection(consts.CacheCollection).Doc(key).Set(firebase.Ctx, map[string]interface{}{
		consts.FieldTimestamp: time.Now(),
		consts.FieldData:      data,
	})
	return err
}

// getCache retrieves data from the cache and returns entry data and timestamp
func getCache[T any](ctx context.Context, collection, key string, maxAge time.Duration) (*T, time.Time, error) {
	doc, err := firebase.Client.Collection(collection).Doc(key).Get(ctx)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf(consts.CacheMissFor, key, err)
	}

	var entry struct {
		Timestamp time.Time `firestore:"timestamp"`
		Data      T         `firestore:"data"`
	}

	if err := doc.DataTo(&entry); err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to decode cache entry for %s: %w", key, err)
	}

	return &entry.Data, entry.Timestamp, nil
}

// GetCachedCountryInfo retrieves cached country information
func GetCachedCountryInfo(ctx context.Context, country string, maxAge time.Duration) (interface{}, bool, error) {
	key := consts.CacheCountryInfoPrefix + country
	data, ts, err := getCache[interface{}](ctx, consts.CacheCollection, key, maxAge)
	if err != nil {
		return nil, false, err
	}
	if isCacheExpired(ts, maxAge) {
		return nil, false, fmt.Errorf(consts.CacheExpiredFor, key)
	}
	return *data, true, nil
}

// SaveCountryInfoToCache stores country information in the cache
func SaveCountryInfoToCache(country string, data interface{}) error {
	key := CountryCacheKey(country)
	return setCache(key, data)
}

// GetCachedWeather retrieves cached weather information
func GetCachedWeather(ctx context.Context, locationKey string, maxAge time.Duration) (interface{}, bool, error) {
	key := consts.CacheWeatherPrefix + locationKey
	data, ts, err := getCache[interface{}](ctx, consts.CacheCollection, key, maxAge)
	if err != nil {
		return nil, false, err
	}
	if isCacheExpired(ts, maxAge) {
		return nil, false, fmt.Errorf(consts.CacheExpiredFor, key)
	}
	return *data, true, nil
}

// SaveWeatherToCache stores weather information in the cache
func SaveWeatherToCache(key string, data interface{}) error {
	return setCache(key, data)
}

// GetCachedCurrencyRates retrieves cached currency rates
func GetCachedCurrencyRates(ctx context.Context, baseCurrency string, maxAge time.Duration) (interface{}, bool, error) {
	key := consts.CacheCurrencyPrefix + baseCurrency
	data, ts, err := getCache[interface{}](ctx, consts.CacheCollection, key, maxAge)
	if err != nil {
		return nil, false, err
	}
	if isCacheExpired(ts, maxAge) {
		return nil, false, fmt.Errorf(consts.CacheExpiredFor, key)
	}
	return *data, true, nil
}

// SaveCurrencyRatesToCache stores currency rates in the cache
func SaveCurrencyRatesToCache(baseCurrency string, data interface{}) error {
	key := CurrencyCacheKey(baseCurrency)
	return setCache(key, data)
}
