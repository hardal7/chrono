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

func Delete(ctx context.Context, r dto.DeleteTopicRequest) error {
	err := db.Queries.DeleteTopic(ctx, query.DeleteTopicParams{
		OwnerID: ctx.Value(middleware.UserID).(uuid.UUID),
		Name:    r.Name,
	})
	if err != nil {
		return fmt.Errorf("Failed to delete topic: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
