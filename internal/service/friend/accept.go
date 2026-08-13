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

func AcceptRequest(ctx context.Context, r dto.AcceptFriendRequestRequest) error {
	logger.Info("Accepting friend request", "sender", r.Username)

	u, err := conn.Queries.GetUserByUsername(ctx, r.Username)
	if err != nil {
		logger.Error("Failed to find user", err)
		return err
	}

	err = conn.Queries.AcceptFriendRequest(ctx, db.AcceptFriendRequestParams{
		RecipientID: ctx.Value(middleware.UserID).(uuid.UUID),
		SenderID:    u.ID,
	})
	if err != nil {
		logger.Error("Failed to accept friend request", err)
		return err
	}
	logger.Info("Accepted friend request", "sender", r.Username)
	return nil
}
