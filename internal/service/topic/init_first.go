package topic

import (
	"context"
	"errors"
	"fmt"

	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/jackc/pgx/v5"
)

const firstTopic string = "General"

func InitFirst(ctx context.Context) error {
	userID := middleware.UserID(ctx)

	_, err := db.Queries.GetTopicByOwnerAndName(ctx, query.GetTopicByOwnerAndNameParams{
		Name:    firstTopic,
		OwnerID: userID,
	})
	if !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("First topic already initialized")
	}

	err = db.Queries.CreateTopic(ctx, query.CreateTopicParams{
		Name:    firstTopic,
		OwnerID: userID,
	})
	if err != nil {
		return fmt.Errorf("Failed to create topic: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
