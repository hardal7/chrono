package runner

import (
	"context"
	"time"

	conn "github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/util/logger"
)

// TODO: Set streaks for topics here
func ResetToday(ctx context.Context) {
	logger.Info("Started runner", "name", "reset_today")
	for {
		now := time.Now()
		nextDay := now.AddDate(0, 0, 1)
		tomorrow := time.Date(
			nextDay.Year(),
			nextDay.Month(),
			nextDay.Day(),
			0, 0, 0, 0,
			time.UTC,
		)

		timer := time.NewTimer(time.Until(tomorrow))
		<-timer.C
		err := conn.Queries.ResetTimeTrackedToday(ctx)
		if err != nil {
			logger.Debug(err.Error())
			logger.Error("Failed to reset time tracked for today", "date", time.Now().String())
		} else {
			logger.Info("Reset times tracked today")
		}
	}
}
