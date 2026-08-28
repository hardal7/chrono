package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5"
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
		err := updateStreaks(ctx)
		if err != nil {
			logger.Error(err.Error())
		}
		err = resetTodayTimes(ctx)
		if err != nil {
			logger.Error(err.Error())
		}
	}
}

func resetTodayTimes(ctx context.Context) error {
	logger.Info("Reseting times tracked today")

	tx, err := db.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("Failed to begin new transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	err = db.Queries.WithTx(tx).ResetTopicTimeTrackedToday(ctx)
	if err != nil {
		return fmt.Errorf("Failed to reset time tracked for today: type=topic, date=%s error=%q", time.Now().String(), err)
	}
	err = db.Queries.WithTx(tx).ResetUserTimeTrackedToday(ctx)
	if err != nil {
		return fmt.Errorf("Failed to reset time tracked for today: type=user, date=%s error=%q", time.Now().String(), err)
	}
	err = db.Queries.WithTx(tx).ResetSessionParticipantTimeTrackedToday(ctx)
	if err != nil {
		return fmt.Errorf("Failed to reset time tracked for today: type=session_participant, date=%s error=%q", time.Now().String(), err)
	}

	logger.Info("Reset times tracked today")
	return nil
}

func updateStreaks(ctx context.Context) error {
	logger.Info("Updating streaks")

	topics, err := db.Queries.GetTopicsAll(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get topics: %q", err)
	}
	for _, topic := range topics {
		if topic.TodayTimeTrackedSeconds != 0 {
			err = db.Queries.IncreaseStreak(ctx, topic.ID)
		} else {
			err = db.Queries.LoseStreak(ctx, topic.ID)
		}

		if err != nil {
			return fmt.Errorf("Failed to update streak: %q", err)
		}
	}

	logger.Info("Updated streaks", "topics", len(topics))
	return nil
}
