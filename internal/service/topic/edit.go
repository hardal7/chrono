package topic

import (
	"context"

	"github.com/google/uuid"
	conn "github.com/hardal7/chrono/internal/db"
	db "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
)

func Edit(ctx context.Context, r dto.EditTopicRequest) error {
	t, err := conn.Queries.GetTopicByOwnerAndName(ctx, db.GetTopicByOwnerAndNameParams{
		Name:    r.Name,
		OwnerID: ctx.Value(middleware.UserID).(uuid.UUID),
	})
	if err != nil {
		return err
	}
	logger.Info("Editing topic", "topicName", t.Name)
	if r.Delete {
		logger.Info("Deleting topic", "topicName", t.Name)
		err := conn.Queries.DeleteTopic(ctx, t.ID)
		if err != nil {
			logger.Error("Failed to delete topic", err)
			return err
		}
		logger.Info("Deleted topic")
		return nil
	}
	if r.NewName == "" {
		logger.Info("Topic not changed")
		return nil
	}

	t.Name = r.NewName
	err = conn.Queries.UpdateTopic(ctx, db.UpdateTopicParams{
		ID:   t.ID,
		Name: t.Name,
	})
	if err != nil {
		logger.Error("Failed to update topic", err)
		return err
	}
	logger.Info("Changed topic name", "newName", r.NewName)
	return nil
}
