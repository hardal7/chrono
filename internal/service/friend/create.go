package friend

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
)

func CreateRequest(ctx context.Context, r dto.CreateFriendRequestRequest) error {
	err := db.Queries.CreateFriendRequest(ctx, query.CreateFriendRequestParams{
		SenderID: ctx.Value(middleware.UserID).(uuid.UUID),
		Username: r.Username,
	})
	if err != nil {
		return fmt.Errorf("Failed to create friend request: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
