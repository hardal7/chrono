package topicevent

import (
	"context"

	"github.com/google/uuid"
	conn "github.com/hardal7/chrono/internal/db"
	db "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
)

func Track(ctx context.Context, r dto.TrackTopicEventRequest) error {
	logger.Info("Tracking time", "topic", r.Topic)
	t, err := conn.Queries.GetTopicOfUserByName(ctx, db.GetTopicOfUserByNameParams{
		Name:            r.Topic,
		CreatedByUserid: ctx.Value(middleware.UserID).(uuid.UUID),
	})
	if err != nil {
		logger.Error("Failed to get topic by username", err)
		return err
	}
	err = conn.Queries.TrackTopicTime(ctx, db.TrackTopicTimeParams{
		ID:                 t.ID,
		CreatedByUserid:    ctx.Value(middleware.UserID).(uuid.UUID),
		TimeTrackedSeconds: int32(r.TimeSeconds),
	})
	if err != nil {
		logger.Error("Failed to track topic time", err)
		return err
	}
	err = conn.Queries.TrackUserTime(ctx, db.TrackUserTimeParams{
		ID:                      ctx.Value(middleware.UserID).(uuid.UUID),
		TotalTimeTrackedSeconds: int32(r.TimeSeconds),
	})
	if err != nil {
		logger.Error("Failed to track user time", err)
		return err
	}
	err = conn.Queries.CreateTopicEvent(ctx, db.CreateTopicEventParams{
		UserID:             ctx.Value(middleware.UserID).(uuid.UUID),
		TopicID:            t.ID,
		TimeTrackedSeconds: int32(r.TimeSeconds),
		Date:               r.Date,
	})
	if err != nil {
		logger.Error("Failed to create topic event", err)
		return err
	}
	logger.Info("Tracked time", "topic", r.Topic)
	return nil
}
