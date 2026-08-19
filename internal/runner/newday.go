package runner

import (
	"context"
	"time"

	conn "github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/util/logger"
)

func nextDate(days int) time.Time {
	now := time.Now()
	nextDay := now.AddDate(0, 0, days)
	return time.Date(
		nextDay.Year(),
		nextDay.Month(),
		nextDay.Day(),
		0, 0, 0, 0,
		time.UTC,
	)
}

func NewDay(ctx context.Context) {
	logger.Info("Started runner", "name", "new_day")
	for {
		timer := time.NewTimer(time.Until(nextDate(1)))
		<-timer.C
		updateStreaks(ctx)
		resetTodayTimes(ctx)
	}
}

func resetTodayTimes(ctx context.Context) {
	logger.Info("Reseting times tracked today")
	err := conn.Queries.ResetTopicTimeTrackedToday(ctx)
	if err != nil {
		logger.Error("Failed to reset time tracked for today", "date", time.Now().String(), "error", err)
		return
	}
	err = conn.Queries.ResetUserTimeTrackedToday(ctx)
	if err != nil {
		logger.Error("Failed to reset time tracked for today", "date", time.Now().String(), "error", err)
		return
	}
	logger.Info("Reset times tracked today")
}

func updateStreaks(ctx context.Context) {
	logger.Info("Updating streaks")
	topics, err := conn.Queries.GetTopicsAll(ctx)
	if err != nil {
		logger.Error("Failed to get topics", err)
		return
	}
	for _, topic := range topics {
		if topic.TodayTimeTrackedSeconds != 0 {
			err = conn.Queries.IncreaseStreak(ctx, topic.ID)
		} else {
			err = conn.Queries.LoseStreak(ctx, topic.ID)
		}

		if err != nil {
			logger.Error("Failed to update streak", err)
		}
	}
	logger.Info("Updated streaks", "topics", len(topics))
}
