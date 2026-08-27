package topic

import (
	"context"
	"fmt"

	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
)

func Edit(ctx context.Context, r dto.EditTopicRequest) error {
	userID := middleware.UserID(ctx)

	t, err := db.Queries.GetTopicByOwnerAndName(ctx, query.GetTopicByOwnerAndNameParams{
		Name:    r.Name,
		OwnerID: userID,
	})
	if err != nil {
		return fmt.Errorf("Failed to get topic: %w: %w", db.ErrRunQuery, err)
	}

	if r.NewName == "" {
		return nil
	}

	t.Name = r.NewName
	err = db.Queries.UpdateTopic(ctx, query.UpdateTopicParams{
		ID:   t.ID,
		Name: t.Name,
	})
	if err != nil {
		return fmt.Errorf("Failed to update topic: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
