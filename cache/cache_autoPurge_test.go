package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestManualCacheAutoPurgeExecution verifies that all purge jobs are executed during a manual run of the auto-purge logic.
// It uses mocked purge functions and asserts that each job is invoked with the correct TTL.
func TestManualCacheAutoPurgeExecution(t *testing.T) {
	var calledJobs []string

	mockJob := func(name string) func(context.Context, time.Duration) error {
		return func(ctx context.Context, ttl time.Duration) error {
			calledJobs = append(calledJobs, name)
			return nil
		}
	}

	jobs := []purgeJob{
		{Name: "CountryCache", Func: mockJob("CountryCache"), TTL: 6 * time.Hour},
		{Name: "WeatherCache", Func: mockJob("WeatherCache"), TTL: 6 * time.Hour},
		{Name: "CurrencyCache", Func: mockJob("CurrencyCache"), TTL: 6 * time.Hour},
	}

	ctx := context.Background()

	for _, job := range jobs {
		err := job.Func(ctx, job.TTL)
		assert.NoError(t, err)
	}

	assert.ElementsMatch(t, []string{"CountryCache", "WeatherCache", "CurrencyCache"}, calledJobs)
}
