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
		logger.Debug(err.Error())
		logger.Error("Failed to reset time tracked for today", "date", time.Now().String())
		return
	}
	err = conn.Queries.ResetUserTimeTrackedToday(ctx)
	if err != nil {
		logger.Debug(err.Error())
		logger.Error("Failed to reset time tracked for today", "date", time.Now().String())
		return
	}
	logger.Info("Reset times tracked today")
}

func updateStreaks(ctx context.Context) {
	logger.Info("Updating streaks")
	topics, err := conn.Queries.GetTopicsAll(ctx)
	if err != nil {
		logger.Debug(err.Error())
		logger.Error("Failed to get topics")
		return
	}
	for _, v := range topics {
		if v.TodayTimeTrackedSeconds != 0 {
			conn.Queries.IncreaseStreak(ctx, v.ID)
		} else {
			conn.Queries.LoseStreak(ctx, v.ID)
		}
	}
	logger.Info("Updated streaks", "topics", len(topics))
}
