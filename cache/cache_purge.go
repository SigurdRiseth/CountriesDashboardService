package cache

import (
	"CountriesDashboardService/consts"
	"CountriesDashboardService/firebase"
	"context"
	"fmt"
	"log"
	"time"
)

// purgeCacheCollection scans a Firestore collection for documents with a timestamp
// older than the given TTL (time-to-live) and deletes them.
// This is a general-purpose internal function used by all specific purge functions.
func purgeCacheCollection(ctx context.Context, collection string, ttl time.Duration) error {
	if !firebase.IsFirebaseClientInitialized() {
		return fmt.Errorf(consts.FTNotInitialized)
	}

	threshold := time.Now().Add(-ttl)
	query := firebase.GetClient().Collection(collection).Where(consts.TimeStamp, consts.LessThanAlligator, threshold)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return fmt.Errorf(consts.FailedToQueryExpiredDocs, err)
	}

	for _, doc := range docs {
		_, err := doc.Ref.Delete(ctx)
		if err != nil {
			log.Printf(consts.FailedToDeleteDocs, doc.Ref.ID, err)
		}
	}

	log.Printf(consts.PurgedExpiredDocs, len(docs), collection)
	return nil
}

// PurgeCountryCache performs a cleanup of outdated country-related cache entries
// by calling the shared purge logic with the appropriate TTL.
func PurgeCountryCache(ctx context.Context, ttl time.Duration) error {
	return purgeCacheCollection(ctx, consts.CacheCollection, ttl) // Assuming all cached in same collection
}

// PurgeWeatherCache removes weather-related cache entries that have expired.
// It uses the general cache purging mechanism with a TTL value.
func PurgeWeatherCache(ctx context.Context, ttl time.Duration) error {
	return purgeCacheCollection(ctx, consts.CacheCollection, ttl)
}

// PurgeCurrencyCache cleans up outdated cached currency rate documents.
// It uses the same collection and TTL-based deletion logic.
func PurgeCurrencyCache(ctx context.Context, ttl time.Duration) error {
	return purgeCacheCollection(ctx, consts.CacheCollection, ttl)
}
