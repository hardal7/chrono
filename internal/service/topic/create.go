package topic

import (
	"context"
	"fmt"

	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/util/requestctx"
)

func Create(ctx context.Context, r dto.CreateTopicRequest) error {
	userID := requestctx.GetUserID(ctx)

	err := db.Queries.CreateTopic(ctx, query.CreateTopicParams{
		Name:    r.Name,
		OwnerID: userID,
	})
	if err != nil {
		return fmt.Errorf("create topic: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
