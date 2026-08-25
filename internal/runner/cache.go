package runner

import (
	"context"
	"time"

	"github.com/hardal7/chrono/internal/util/logger"
)

func nextMinute() time.Time {
	now := time.Now()
	return now.Truncate(time.Minute).Add(time.Minute)
}

func Cache(ctx context.Context) {
	logger.Info("Started runner", "name", "cache")
	for {
		timer := time.NewTimer(time.Until(nextMinute()))
		<-timer.C
		updateCache(ctx)
	}
}

func updateCache(ctx context.Context) {

}
