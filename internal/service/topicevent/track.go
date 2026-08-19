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
	logger.Debug("Tracking time", "topic", r.Topic)
	topic, err := conn.Queries.GetTopicByOwnerAndName(ctx, db.GetTopicByOwnerAndNameParams{
		Name:    r.Topic,
		OwnerID: ctx.Value(middleware.UserID).(uuid.UUID),
	})
	if err != nil {
		logger.Debug("Failed to get topic by username", err)
		return err
	}

	err = conn.Queries.TrackTopicTime(ctx, db.TrackTopicTimeParams{
		ID:          topic.ID,
		OwnerID:     ctx.Value(middleware.UserID).(uuid.UUID),
		TimeTracked: int32(r.TimeSeconds),
	})
	if err != nil {
		logger.Debug("Failed to track topic time", err)
		return err
	}

	err = conn.Queries.TrackUserTime(ctx, db.TrackUserTimeParams{
		ID:          ctx.Value(middleware.UserID).(uuid.UUID),
		TimeTracked: int32(r.TimeSeconds),
	})
	if err != nil {
		logger.Debug("Failed to track user time", err)
		return err
	}

	err = conn.Queries.CreateTopicEvent(ctx, db.CreateTopicEventParams{
		UserID:             ctx.Value(middleware.UserID).(uuid.UUID),
		TopicID:            topic.ID,
		TimeTrackedSeconds: int32(r.TimeSeconds),
		CreatedAt:          r.Date,
	})
	if err != nil {
		logger.Debug("Failed to create topic event", err)
		return err
	}

	logger.Debug("Tracked time", "topic", r.Topic)
	return nil
}
