package topic

import (
	"context"
	"fmt"

	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/util/requestctx"
)

const firstTopic string = "General"

func InitFirst(ctx context.Context) error {
	userID := requestctx.GetUserID(ctx)

	err := db.Queries.CreateTopic(ctx, query.CreateTopicParams{
		Name:    firstTopic,
		OwnerID: userID,
	})
	if err != nil {
		return fmt.Errorf("create topic: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
