package friend

import (
	"context"

	"github.com/google/uuid"
	conn "github.com/hardal7/chrono/internal/db"
	db "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
)

func Remove(ctx context.Context, r dto.RemoveFriendRequest) error {
	logger.Info("Removing friend", "friend", r.Username)

	u, err := conn.Queries.GetUserByUsername(ctx, r.Username)
	if err != nil {
		logger.Error("Failed to find friend", err)
		return err
	}

	err = conn.Queries.DeleteFriend(ctx, db.DeleteFriendParams{
		SenderID:    ctx.Value(middleware.UserID).(uuid.UUID),
		RecipientID: u.ID,
	})
	if err != nil {
		logger.Error("Failed to remove friend", err)
		return err
	}
	logger.Info("Removed friend", "friend", r.Username)
	return nil
}
