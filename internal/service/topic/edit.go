package topic

import (
	"context"
	"fmt"

	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/auth"
)

func Edit(ctx context.Context, r dto.EditTopicRequest) error {
	userID := auth.UserID(ctx)

	err := db.Queries.UpdateTopic(ctx, query.UpdateTopicParams{
		OwnerID: userID,
		Name:    r.Name,
		NewName: r.NewName,
	})
	if err != nil {
		return fmt.Errorf("Failed to update topic: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
