package cache

import (
	"CountriesDashboardService/consts"
	"CountriesDashboardService/firebase"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestPurgeCountryCache checks that expired country cache documents are deleted while fresh documents remain.
func TestPurgeCountryCache(t *testing.T) {
	ctx := context.Background()

	// Write 1 fresh doc and 1 expired doc
	freshID := "fresh-doc"
	expiredID := "expired-doc"
	now := time.Now()

	// Fresh doc
	_, err := firebase.GetClient().Collection(consts.CacheCollection).Doc(freshID).Set(ctx, map[string]interface{}{
		consts.FieldTimestamp: now,
		consts.FieldData:      "still valid",
	})
	assert.NoError(t, err)

	// Expired doc
	_, err = firebase.GetClient().Collection(consts.CacheCollection).Doc(expiredID).Set(ctx, map[string]interface{}{
		consts.FieldTimestamp: now.Add(-10 * time.Second),
		consts.FieldData:      "old data",
	})
	assert.NoError(t, err)

	// Purge with TTL = 1 hour
	err = PurgeCountryCache(ctx, 5*time.Second)
	assert.NoError(t, err)

	// Check that fresh doc still exists
	_, err = firebase.GetClient().Collection(consts.CacheCollection).Doc(freshID).Get(ctx)
	assert.NoError(t, err, "fresh doc should not be deleted")

	// Check that expired doc was deleted
	_, err = firebase.GetClient().Collection(consts.CacheCollection).Doc(expiredID).Get(ctx)
	assert.Error(t, err, "expired doc should be deleted")
}

// TestPurgeWeatherCache checks that expired weather cache documents are deleted while fresh documents remain.
func TestPurgeWeatherCache(t *testing.T) {
	ctx := context.Background()

	// Write 1 fresh doc and 1 expired doc
	freshID := "fresh-doc"
	expiredID := "expired-doc"
	now := time.Now()

	// Fresh doc
	_, err := firebase.GetClient().Collection(consts.CacheCollection).Doc(freshID).Set(ctx, map[string]interface{}{
		consts.FieldTimestamp: now,
		consts.FieldData:      "still valid",
	})
	assert.NoError(t, err)

	// Expired doc
	_, err = firebase.GetClient().Collection(consts.CacheCollection).Doc(expiredID).Set(ctx, map[string]interface{}{
		consts.FieldTimestamp: now.Add(-2 * time.Hour),
		consts.FieldData:      "old data",
	})
	assert.NoError(t, err)

	// Purge with TTL = 1 hour
	err = PurgeWeatherCache(ctx, 1*time.Hour)
	assert.NoError(t, err)

	// Check that fresh doc still exists
	_, err = firebase.GetClient().Collection(consts.CacheCollection).Doc(freshID).Get(ctx)
	assert.NoError(t, err, "fresh doc should not be deleted")

	// Check that expired doc was deleted
	_, err = firebase.GetClient().Collection(consts.CacheCollection).Doc(expiredID).Get(ctx)
	assert.Error(t, err, "expired doc should be deleted")
}

// TestPurgeCurrencyCache checks that expired currency cache documents are deleted while fresh documents remain.
func TestPurgeCurrencyCache(t *testing.T) {
	ctx := context.Background()

	// Write 1 fresh doc and 1 expired doc
	freshID := "fresh-doc"
	expiredID := "expired-doc"
	now := time.Now()

	// Fresh doc
	_, err := firebase.GetClient().Collection(consts.CacheCollection).Doc(freshID).Set(ctx, map[string]interface{}{
		consts.FieldTimestamp: now,
		consts.FieldData:      "still valid",
	})
	assert.NoError(t, err)

	// Expired doc
	_, err = firebase.GetClient().Collection(consts.CacheCollection).Doc(expiredID).Set(ctx, map[string]interface{}{
		consts.FieldTimestamp: now.Add(-2 * time.Hour),
		consts.FieldData:      "old data",
	})
	assert.NoError(t, err)

	// Purge with TTL = 1 hour
	err = PurgeCurrencyCache(ctx, 1*time.Hour)
	assert.NoError(t, err)

	// Check that fresh doc still exists
	_, err = firebase.GetClient().Collection(consts.CacheCollection).Doc(freshID).Get(ctx)
	assert.NoError(t, err, "fresh doc should not be deleted")

	// Check that expired doc was deleted
	_, err = firebase.GetClient().Collection(consts.CacheCollection).Doc(expiredID).Get(ctx)
	assert.Error(t, err, "expired doc should be deleted")
}
