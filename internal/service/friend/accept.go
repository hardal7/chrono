package friend

import (
	"context"
	"fmt"

	"github.com/hardal7/chrono/internal/auth"
	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
)

func AcceptRequest(ctx context.Context, r dto.AcceptFriendRequestRequest) error {
	userID := auth.UserID(ctx)

	err := db.Queries.AcceptFriendRequest(ctx, query.AcceptFriendRequestParams{
		RecipientID: userID,
		Username:    r.Username,
	})
	if err != nil {
		return fmt.Errorf("Failed to accept friend request: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
