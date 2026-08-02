package topicevent

import (
	"context"
	"time"

	conn "github.com/hardal7/chrono/internal/db"
	db "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5/pgtype"
)

func Track(ctx context.Context, r dto.TrackTopicEventRequest) error {
	logger.Info("Tracking time", "topic", r.Topic)
	t, err := conn.Queries.GetTopicOfUserByName(ctx, db.GetTopicOfUserByNameParams{
		Name:            r.Topic,
		CreatedByUserid: ctx.Value(middleware.UserID).(int32),
	})
	if err != nil {
		return err
	}
	err = conn.Queries.TrackTopicTime(ctx, db.TrackTopicTimeParams{
		ID:                 t.ID,
		CreatedByUserid:    ctx.Value(middleware.UserID).(int32),
		TimeTrackedSeconds: int32(r.TimeSeconds),
	})
	if err != nil {
		return err
	}
	err = conn.Queries.TrackUserTime(ctx, db.TrackUserTimeParams{
		ID:                 t.ID,
		TimeTrackedSeconds: int32(r.TimeSeconds),
	})
	if err != nil {
		return err
	}
	err = conn.Queries.CreateTopicEvent(ctx, db.CreateTopicEventParams{
		UserID:             ctx.Value(middleware.UserID).(int32),
		TopicID:            t.ID,
		TimeTrackedSeconds: int32(r.TimeSeconds),
		Date:               pgtype.Timestamptz{Time: time.Now()},
	})
	if err != nil {
		return err
	}
	logger.Info("Tracked time", "topic", r.Topic)
	return nil
}
