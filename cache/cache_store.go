package cache

import (
	"CountriesDashboardService/consts"
	"CountriesDashboardService/firebase"
	"context"
	"fmt"
	"time"
)

// CacheItem represents a single cached document in Firestore.
// It contains a timestamp used for TTL logic and the associated data payload.
type CacheItem struct {
	Timestamp time.Time   `firestore:"timestamp"`
	Data      interface{} `firestore:"data"`
}

// isCacheExpired checks whether the provided timestamp is older than the given TTL (maxAge).
// Returns true if the cached item should be considered expired.
func isCacheExpired(ts time.Time, maxAge time.Duration) bool {
	return time.Since(ts) > maxAge
}

// setCache stores a new entry in the Firestore cache collection with the current timestamp.
// It overwrites any existing document with the same key.
func setCache(key string, data interface{}) error {
	if !firebase.IsFirebaseClientInitialized() {
		return fmt.Errorf(consts.FTNotInitialized)
	}
	_, err := firebase.GetClient().Collection(consts.CacheCollection).Doc(key).Set(firebase.GetCtx(), map[string]interface{}{
		consts.FieldTimestamp: time.Now(),
		consts.FieldData:      data,
	})
	return err
}

// getCache retrieves a generic document from Firestore, extracting both data and its timestamp.
// It returns a typed pointer to the data, the timestamp, and an error if applicable.
func getCache[T any](ctx context.Context, collection, key string, maxAge time.Duration) (*T, time.Time, error) {
	doc, err := firebase.GetClient().Collection(collection).Doc(key).Get(ctx)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf(consts.CacheMissFor, key, err)
	}

	var entry struct {
		Timestamp time.Time `firestore:"timestamp"`
		Data      T         `firestore:"data"`
	}

	if err := doc.DataTo(&entry); err != nil {
		return nil, time.Time{}, fmt.Errorf(consts.FailedDecodeCacheEntry, key, err)
	}

	return &entry.Data, entry.Timestamp, nil
}

// GetCachedCountryInfo attempts to retrieve country data from the cache for the given country name or code.
// Returns the cached data, a boolean indicating if it was found and valid, and an error if applicable.
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

// SaveCountryInfoToCache stores country-specific data in the cache using a generated cache key.
func SaveCountryInfoToCache(country string, data interface{}) error {
	key := CountryCacheKey(country)
	return setCache(key, data)
}

// GetCachedWeather retrieves weather data from the cache for a specific location key.
// Returns the cached data, a found flag, and any error encountered.
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

// SaveWeatherToCache writes weather data to the cache for the given key.
func SaveWeatherToCache(key string, data interface{}) error {
	return setCache(key, data)
}

// GetCachedCurrencyRates retrieves cached currency exchange rate data for a base currency.
// Returns the cached data, a valid flag, and any error encountered.
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

// SaveCurrencyRatesToCache writes currency exchange rate data into the cache using the base currency as key.
func SaveCurrencyRatesToCache(baseCurrency string, data interface{}) error {
	key := CurrencyCacheKey(baseCurrency)
	return setCache(key, data)
}
