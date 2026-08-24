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

func Remove(ctx context.Context, r dto.RemoveFriendRequest) error {
	err := db.Queries.DeleteFriend(ctx, query.DeleteFriendParams{
		SenderID: ctx.Value(middleware.UserID).(uuid.UUID),
		Username: r.Username,
	})
	if err != nil {
		return fmt.Errorf("Failed to remove friend: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
