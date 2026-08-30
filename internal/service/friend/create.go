package friend

import (
	"context"
	"fmt"

	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/auth"
)

func CreateRequest(ctx context.Context, r dto.CreateFriendRequestRequest) error {
	userID := auth.UserID(ctx)

	err := db.Queries.CreateFriendRequest(ctx, query.CreateFriendRequestParams{
		SenderID: userID,
		Username: r.Username,
	})
	if err != nil {
		return fmt.Errorf("Failed to create friend request: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
