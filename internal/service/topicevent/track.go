package topicevent

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
)

func Track(ctx context.Context, r dto.TrackTopicEventRequest) error {
	topic, err := db.Queries.GetTopicByOwnerAndName(ctx, query.GetTopicByOwnerAndNameParams{
		Name:    r.Topic,
		OwnerID: ctx.Value(middleware.UserID).(uuid.UUID),
	})
	if err != nil {
		return fmt.Errorf("Failed to get topic by username: %w: %w", db.ErrRunQuery, err)
	}

	err = db.Queries.TrackTopicTime(ctx, query.TrackTopicTimeParams{
		ID:          topic.ID,
		OwnerID:     ctx.Value(middleware.UserID).(uuid.UUID),
		TimeTracked: int32(r.TimeSeconds),
	})
	if err != nil {
		return fmt.Errorf("Failed to track topic time: %w: %w", db.ErrRunQuery, err)
	}

	err = db.Queries.TrackUserTime(ctx, query.TrackUserTimeParams{
		ID:          ctx.Value(middleware.UserID).(uuid.UUID),
		TimeTracked: int32(r.TimeSeconds),
	})
	if err != nil {
		return fmt.Errorf("Failed to track user time: %w: %w", db.ErrRunQuery, err)
	}

	err = db.Queries.CreateTopicEvent(ctx, query.CreateTopicEventParams{
		UserID:             ctx.Value(middleware.UserID).(uuid.UUID),
		TopicID:            topic.ID,
		TimeTrackedSeconds: int32(r.TimeSeconds),
		CreatedAt:          r.Date,
	})
	if err != nil {
		return fmt.Errorf("Failed to create topic event: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
