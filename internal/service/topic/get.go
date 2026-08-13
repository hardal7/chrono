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

func Get(ctx context.Context, r dto.GetTopicRequest) (dto.GetTopicResponse, error) {
	logger.Info("Getting topic", "name", r.Name)
	t, err := conn.Queries.GetTopicByOwnerAndName(ctx, db.GetTopicByOwnerAndNameParams{
		Name:    r.Name,
		OwnerID: ctx.Value(middleware.UserID).(uuid.UUID),
	})
	if err != nil {
		logger.Error("Failed to get topic by username", err)
		return dto.GetTopicResponse{}, err
	}
	resp := dto.GetTopicResponse{
		TotalTime: int(t.TotalTimeTrackedSeconds),
		Streak:    int(t.Streak),
	}
	logger.Info("Got topic", "name", r.Name)
	return resp, nil
}
