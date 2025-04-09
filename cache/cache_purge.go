package cache

import (
	"CountriesDashboardService/consts"
	"CountriesDashboardService/firebase"
	"context"
	"fmt"
	"log"
	"time"
)

// purgeCacheCollection removes documents older than a specified TTL from a Firestore collection.
func purgeCacheCollection(ctx context.Context, collection string, ttl time.Duration) error {
	if firebase.Client == nil {
		return fmt.Errorf(consts.FTNotInitialized)
	}

	threshold := time.Now().Add(-ttl)
	query := firebase.Client.Collection(collection).Where("timestamp", "<", threshold)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return fmt.Errorf("failed to query expired documents: %w", err)
	}

	for _, doc := range docs {
		_, err := doc.Ref.Delete(ctx)
		if err != nil {
			log.Printf("Failed to delete doc %s: %v", doc.Ref.ID, err)
		}
	}

	log.Printf("Purged %d expired documents from %s", len(docs), collection)
	return nil
}

// PurgeCountryCache deletes outdated country cache entries.
func PurgeCountryCache(ctx context.Context, ttl time.Duration) error {
	return purgeCacheCollection(ctx, consts.CacheCollection, ttl) // Assuming all cached in same collection
}

// PurgeWeatherCache deletes outdated weather cache entries.
func PurgeWeatherCache(ctx context.Context, ttl time.Duration) error {
	return purgeCacheCollection(ctx, consts.CacheCollection, ttl)
}

// PurgeCurrencyCache deletes outdated currency cache entries.
func PurgeCurrencyCache(ctx context.Context, ttl time.Duration) error {
	return purgeCacheCollection(ctx, consts.CacheCollection, ttl)
}
