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

func GetNamed(ctx context.Context, r dto.GetTopicNamedRequest) (dto.GetTopicNamedResponse, error) {
	logger.Debug("Getting topic", "topicName", r.Name)
	t, err := conn.Queries.GetTopicByOwnerAndName(ctx, db.GetTopicByOwnerAndNameParams{
		Name:    r.Name,
		OwnerID: ctx.Value(middleware.UserID).(uuid.UUID),
	})
	if err != nil {
		logger.Debug("Failed to get topic by username", err)
		return dto.GetTopicNamedResponse{}, err
	}
	resp := dto.GetTopicNamedResponse{
		TotalTime: int(t.TotalTimeTrackedSeconds),
		Streak:    int(t.Streak),
	}
	logger.Debug("Got topic", "topicName", r.Name)
	return resp, nil
}
