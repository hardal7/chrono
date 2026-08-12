package topic

import (
	"context"
	"errors"

	"github.com/google/uuid"
	conn "github.com/hardal7/chrono/internal/db"
	db "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5"
)

func Create(ctx context.Context, r dto.CreateTopicRequest) error {
	logger.Info("Creating topic", "topicName", r.Name)

	_, err := conn.Queries.GetTopicByOwnerAndName(ctx, db.GetTopicByOwnerAndNameParams{
		OwnerID: ctx.Value(middleware.UserID).(uuid.UUID),
		Name:    r.Name,
	})
	if err != pgx.ErrNoRows {
		if err == nil {
			logger.Error("Topic with name exists")
			return errors.New("topic with name exists")
		} else {
			logger.Error("Failed to check if topic is duplicate", err)
			return err
		}
	}

	err = conn.Queries.CreateTopic(ctx, db.CreateTopicParams{
		Name:    r.Name,
		OwnerID: ctx.Value(middleware.UserID).(uuid.UUID),
	})
	if err != nil {
		logger.Error("Failed to create topic", err)
		return err
	}
	logger.Info("Created topic", "topicName", r.Name)
	return nil
}
