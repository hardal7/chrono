package topic

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
)

func GetNamed(ctx context.Context, r dto.GetTopicNamedRequest) (dto.GetTopicNamedResponse, error) {
	resp := dto.GetTopicNamedResponse{}

	t, err := db.Queries.GetTopicByOwnerAndName(ctx, query.GetTopicByOwnerAndNameParams{
		Name:    r.Name,
		OwnerID: ctx.Value(middleware.UserID).(uuid.UUID),
	})
	if err != nil {
		return resp, fmt.Errorf("Failed to get topic by username: %w: %w", db.ErrRunQuery, err)
	}

	resp = dto.GetTopicNamedResponse{
		TotalTime: int(t.TotalTimeTrackedSeconds),
		Streak:    int(t.Streak),
	}

	return resp, nil
}
