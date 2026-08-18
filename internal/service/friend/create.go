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

func CreateRequest(ctx context.Context, r dto.CreateFriendRequestRequest) error {
	logger.Info("Creating friend request", "recipient", r.Username)

	err := conn.Queries.CreateFriendRequest(ctx, db.CreateFriendRequestParams{
		SenderID: ctx.Value(middleware.UserID).(uuid.UUID),
		Username: r.Username,
	})
	if err != nil {
		logger.Error("Failed to create friend request", err)
		return err
	}
	logger.Info("Created friend request", "recipient", r.Username)
	return nil
}
