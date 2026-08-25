package runner

import (
	"context"
	"time"

	conn "github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/util/logger"
)

func nextMidnight() time.Time {
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)

	return time.Date(
		tomorrow.Year(),
		tomorrow.Month(),
		tomorrow.Day(),
		0, 0, 0, 0,
		now.Location(),
	)
}

func NewDay(ctx context.Context) {
	logger.Info("Started runner", "name", "new_day")
	for {
		timer := time.NewTimer(time.Until(nextMidnight()))
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
		logger.Error("Failed to get topics", "error", err)
		return
	}
	for _, topic := range topics {
		if topic.TodayTimeTrackedSeconds != 0 {
			err = conn.Queries.IncreaseStreak(ctx, topic.ID)
		} else {
			err = conn.Queries.LoseStreak(ctx, topic.ID)
		}

		if err != nil {
			logger.Error("Failed to update streak", "error", err)
		}
	}
	logger.Info("Updated streaks", "topics", len(topics))
}
