package topic

import (
	"context"
	"fmt"

	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/util/requestctx"
)

func GetNamed(ctx context.Context, r dto.GetTopicNamedRequest) (dto.GetTopicNamedResponse, error) {
	userID := requestctx.GetUserID(ctx)
	resp := dto.GetTopicNamedResponse{}

	t, err := db.Queries.GetTopicByOwnerAndName(ctx, query.GetTopicByOwnerAndNameParams{
		Name:    r.Name,
		OwnerID: userID,
	})
	if err != nil {
		return resp, fmt.Errorf("get topic by username: %w: %w", db.ErrRunQuery, err)
	}

	resp = dto.GetTopicNamedResponse{
		TotalTime: int(t.TotalTimeTrackedSeconds),
	}

	return resp, nil
}
