package cache

import (
	"CountriesDashboardService/consts"
	"context"
	"log"
	"time"
)

type purgeJob struct {
	Name string
	Func func(context.Context, time.Duration) error
	TTL  time.Duration
}

// StartCacheAutoPurge starts a background task to purge expired cache at intervals.
func StartCacheAutoPurge() {
	ticker := time.NewTicker(6 * time.Hour) // Customize this interval if needed
	go func() {
		for {
			<-ticker.C
			log.Println(consts.AutoPurgedTriggered)

			ctx := context.Background()
			jobs := []purgeJob{
				{Name: consts.CountryCache, Func: PurgeCountryCache, TTL: 6 * time.Hour},
				{Name: consts.WeatherCache, Func: PurgeWeatherCache, TTL: 6 * time.Hour},
				{Name: consts.CurrencyCache, Func: PurgeCurrencyCache, TTL: 6 * time.Hour},
			}

			for _, job := range jobs {
				if err := job.Func(ctx, job.TTL); err != nil {
					log.Printf(consts.ErrorPurging, job.Name, err)
				} else {
					log.Printf(consts.PurgedSuccessfully, job.Name)
				}
			}
		}
	}()
}
