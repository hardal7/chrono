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
	logger.Debug("Accepting friend request", "sender", r.Username)

	err := conn.Queries.AcceptFriendRequest(ctx, db.AcceptFriendRequestParams{
		RecipientID: ctx.Value(middleware.UserID).(uuid.UUID),
		Username:    r.Username,
	})
	if err != nil {
		logger.Debug("Failed to accept friend request", err)
		return err
	}
	logger.Debug("Accepted friend request", "sender", r.Username)
	return nil
}
