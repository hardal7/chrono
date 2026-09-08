package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/util/logger"
)

const (
	nextWeek = 7
	maxRank  = 100
)

func NewWeek(ctx context.Context) {
	logger.Info("Started runner", "name", "new_week")
	for {
		timer := time.NewTimer(time.Until(retrieveDate(nextWeek)))
		<-timer.C
		err := updateLeaderboard(ctx)
		if err != nil {
			logger.Error(err.Error())
		}
		err = resetWeekTimes(ctx)
		if err != nil {
			logger.Error(err.Error())
		}
	}
}

func resetWeekTimes(ctx context.Context) error {
	logger.Info("Reseting times tracked this week")

	err := db.Queries.ResetUserTimeTrackedWeek(ctx)
	if err != nil {
		return fmt.Errorf("Failed to reset time tracked for this week: %w: %w", db.ErrRunQuery, err)
	}

	logger.Info("Reset times tracked this week")
	return nil
}
