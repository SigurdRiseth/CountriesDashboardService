package cache

import (
	"CountriesDashboardService/consts"
	"CountriesDashboardService/firebase"
	"context"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
)

func init() {
	_ = os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8081")
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "test-project")
	if err != nil {
		panic(err)
	}
	firebase.Client = client
	firebase.Ctx = ctx
}

func TestIsCacheExpired(t *testing.T) {
	now := time.Now()
	maxAge := 10 * time.Minute

	assert.False(t, isCacheExpired(now, maxAge), "Fresh cache should not be expired")

	old := now.Add(-15 * time.Minute)
	assert.True(t, isCacheExpired(old, maxAge), "Old cache should be expired")
}

func TestSaveAndGetCountryCache(t *testing.T) {
	ctx := context.Background()
	country := "Testland"
	expected := map[string]interface{}{
		"name": "Testland",
		"data": "interesting",
	}

	err := SaveCountryInfoToCache(country, expected)
	assert.NoError(t, err)

	data, found, err := GetCachedCountryInfo(ctx, country, 10*time.Minute)
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, expected, data)
}

func TestGetCountryCacheExpired(t *testing.T) {
	ctx := context.Background()
	country := "Mongoland"
	expected := map[string]interface{}{
		"name": "Mongoland",
	}

	// Save with past timestamp manually
	key := consts.CacheCountryInfoPrefix + country
	_, err := firebase.Client.Collection(consts.CacheCollection).Doc(key).Set(ctx, map[string]interface{}{
		consts.FieldTimestamp: time.Now().Add(-1 * time.Hour),
		consts.FieldData:      expected,
	})
	assert.NoError(t, err)

	_, found, err := GetCachedCountryInfo(ctx, country, 5*time.Minute)
	assert.Error(t, err)
	assert.False(t, found)
}

func TestSaveAndGetWeatherCache(t *testing.T) {
	ctx := context.Background()
	locationKey := "Oslo"
	expected := map[string]interface{}{
		"temperature": int64(-5),
		"humidity":    int64(75),
	}

	err := SaveWeatherToCache(consts.CacheWeatherPrefix+locationKey, expected)
	assert.NoError(t, err)

	data, found, err := GetCachedWeather(ctx, locationKey, 10*time.Minute)
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, expected, data)
}

func TestGetWeatherCacheExpired(t *testing.T) {
	ctx := context.Background()
	locationKey := "Bergen"
	expected := map[string]interface{}{
		"temperature": 2,
	}

	key := consts.CacheWeatherPrefix + locationKey
	_, err := firebase.Client.Collection(consts.CacheCollection).Doc(key).Set(ctx, map[string]interface{}{
		consts.FieldTimestamp: time.Now().Add(-2 * time.Hour),
		consts.FieldData:      expected,
	})
	assert.NoError(t, err)

	_, found, err := GetCachedWeather(ctx, locationKey, 30*time.Minute)
	assert.Error(t, err)
	assert.False(t, found)
}

func TestSaveAndGetCurrencyRatesCache(t *testing.T) {
	ctx := context.Background()
	baseCurrency := "NOK"
	expected := map[string]interface{}{
		"EUR": 0.088,
		"USD": 0.094,
	}

	err := SaveCurrencyRatesToCache(baseCurrency, expected)
	assert.NoError(t, err)

	data, found, err := GetCachedCurrencyRates(ctx, baseCurrency, 10*time.Minute)
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, expected, data)
}

func TestGetCurrencyRatesCacheExpired(t *testing.T) {
	ctx := context.Background()
	baseCurrency := "SEK"
	expected := map[string]interface{}{
		"EUR": 0.1,
	}

	key := consts.CacheCurrencyPrefix + baseCurrency
	_, err := firebase.Client.Collection(consts.CacheCollection).Doc(key).Set(ctx, map[string]interface{}{
		consts.FieldTimestamp: time.Now().Add(-3 * time.Hour),
		consts.FieldData:      expected,
	})
	assert.NoError(t, err)

	_, found, err := GetCachedCurrencyRates(ctx, baseCurrency, 1*time.Hour)
	assert.Error(t, err)
	assert.False(t, found)
}
