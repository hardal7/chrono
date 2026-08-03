package topic

import (
	"context"

	conn "github.com/hardal7/chrono/internal/db"
	db "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
)

func Get(ctx context.Context, r dto.GetTopicRequest) (dto.GetTopicResponse, error) {
	logger.Info("Getting topic", "name", r.Name)
	t, err := conn.Queries.GetTopicOfUserByName(ctx, db.GetTopicOfUserByNameParams{
		Name:            r.Name,
		CreatedByUserid: ctx.Value(middleware.UserID).(int32),
	})
	if err != nil {
		logger.Error("Failed to get topic by username", err)
		return dto.GetTopicResponse{}, err
	}
	resp := dto.GetTopicResponse{
		TotalTime: int(t.TimeTrackedSeconds),
	}
	logger.Info("Got topic", "name", r.Name)
	return resp, nil
}
