package firebase

import (
	"CountriesDashboardService/consts"
	"fmt"
	"google.golang.org/api/iterator"
	"time"
)

func GetDashboardFromCache(id string) (map[string]interface{}, error) {
	doc, err := Client.Collection(consts.DashboardCacheCollection).Doc(id).Get(Ctx)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = doc.DataTo(&result)
	if err != nil {
		return nil, err
	}

	cachedAtStr, ok := result["cachedAt"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid cache timestamp")
	}
	cachedAt, err := time.Parse(time.RFC3339, cachedAtStr)
	if err != nil || time.Since(cachedAt) > time.Hour {
		return nil, fmt.Errorf("cache expired")
	}

	return result["data"].(map[string]interface{}), nil
}

func StoreDashboardInCache(id string, dashboard map[string]interface{}) error {
	cacheDoc := map[string]interface{}{
		"data":     dashboard,
		"cachedAt": time.Now().Format(time.RFC3339),
	}
	_, err := Client.Collection(consts.DashboardCacheCollection).Doc(id).Set(Ctx, cacheDoc)
	return err
}

func PurgeOldDashboardCache(maxAge time.Duration) error {
	cutoff := time.Now().Add(-maxAge)
	iter := Client.Collection(consts.DashboardCacheCollection).Documents(Ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		cachedAtStr, _ := doc.Data()["cachedAt"].(string)
		cachedAt, err := time.Parse(time.RFC3339, cachedAtStr)
		if err != nil || cachedAt.Before(cutoff) {
			_, _ = doc.Ref.Delete(Ctx)
		}
	}
	return nil
}
