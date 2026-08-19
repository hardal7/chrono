package topic

import (
	"context"

	"github.com/google/uuid"
	conn "github.com/hardal7/chrono/internal/db"
	db "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5"
)

const firstTopic string = "General"

func InitFirst(ctx context.Context) error {
	_, err := conn.Queries.GetTopicByOwnerAndName(ctx, db.GetTopicByOwnerAndNameParams{
		Name:    firstTopic,
		OwnerID: ctx.Value(middleware.UserID).(uuid.UUID),
	})
	if err != pgx.ErrNoRows {
		logger.Debug("First topic already initialized")
		return err
	}

	err = conn.Queries.CreateTopic(ctx, db.CreateTopicParams{
		Name:    firstTopic,
		OwnerID: ctx.Value(middleware.UserID).(uuid.UUID),
	})
	if err != nil {
		return err
	}
	logger.Debug("Initialized first topic")
	return nil
}
