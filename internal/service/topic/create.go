package topic

import (
	"context"
	"errors"

	conn "github.com/hardal7/chrono/internal/db"
	db "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5"
)

func Create(ctx context.Context, r dto.CreateTopicRequest) error {
	logger.Info("Creating topic", "topicName", r.Name)

	_, err := conn.Queries.GetTopicOfUserByName(ctx, db.GetTopicOfUserByNameParams{})
	if err != pgx.ErrNoRows {
		if err == nil {
			logger.Error("Topic with name exists")
			return errors.New("topic with name exists")
		} else {
			logger.Error("Failed to find topic by username", err)
			return err
		}
	}

	err = conn.Queries.CreateTopic(ctx, db.CreateTopicParams{
		Name:               r.Name,
		TimeTrackedSeconds: 0,
		CreatedByUserid:    ctx.Value(middleware.UserID).(int32),
	})
	if err != nil {
		logger.Error("Failed to create topic", err)
		return err
	}
	logger.Info("Created topic", "topicName", r.Name)
	return nil
}
