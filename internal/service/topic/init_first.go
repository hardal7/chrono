package topic

import (
	"context"

	"github.com/google/uuid"
	conn "github.com/hardal7/chrono/internal/db"
	db "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
)

const firstTopic string = "General"

func InitFirst(ctx context.Context) error {
	err := conn.Queries.CreateTopic(ctx, db.CreateTopicParams{
		Name:               firstTopic,
		TimeTrackedSeconds: 0,
		CreatedByUserid:    ctx.Value(middleware.UserID).(uuid.UUID),
	})
	if err != nil {
		logger.Warn("Failed to initialize first topic")
		return err
	}
	logger.Info("Initialized first topic")
	return nil
}
