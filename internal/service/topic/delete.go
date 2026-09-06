package topic

import (
	"context"
	"fmt"

	"github.com/hardal7/chrono/internal/auth"
	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
)

func Delete(ctx context.Context, r dto.DeleteTopicRequest) error {
	userID := auth.UserID(ctx)

	err := db.Queries.DeleteTopic(ctx, query.DeleteTopicParams{
		OwnerID: userID,
		Name:    r.Name,
	})
	if err != nil {
		return fmt.Errorf("Failed to delete topic: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
