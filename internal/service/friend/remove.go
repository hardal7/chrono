package friend

import (
	"context"
	"fmt"

	"github.com/hardal7/chrono/internal/auth"
	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
)

func Remove(ctx context.Context, r dto.RemoveFriendRequest) error {
	userID := auth.UserID(ctx)

	err := db.Queries.DeleteFriend(ctx, query.DeleteFriendParams{
		SenderID: userID,
		Username: r.Username,
	})
	if err != nil {
		return fmt.Errorf("Failed to remove friend: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
