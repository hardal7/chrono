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

func Create(ctx context.Context, r dto.CreateTopicRequest) error {
	err := db.Queries.CreateTopic(ctx, query.CreateTopicParams{
		Name:    r.Name,
		OwnerID: ctx.Value(middleware.UserID).(uuid.UUID),
	})
	if err != nil {
		return fmt.Errorf("Failed to create topic: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
