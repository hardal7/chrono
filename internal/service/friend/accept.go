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

func AcceptRequest(ctx context.Context, r dto.AcceptFriendRequestRequest) error {
	err := db.Queries.AcceptFriendRequest(ctx, query.AcceptFriendRequestParams{
		RecipientID: ctx.Value(middleware.UserID).(uuid.UUID),
		Username:    r.Username,
	})
	if err != nil {
		return fmt.Errorf("Failed to accept friend request: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
