package friend

import (
	"context"

	"github.com/google/uuid"
	conn "github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
)

func GetAll(ctx context.Context) (dto.GetFriendRequestsAllResponse, error) {
	logger.Debug("Getting all friend requests")
	userID := ctx.Value(middleware.UserID).(uuid.UUID)
	r, err := conn.Queries.GetFriendRequests(ctx, userID)
	if err != nil {
		logger.Debug("Failed to get all friend requests", err)
		return dto.GetFriendRequestsAllResponse{}, err
	}

	var requests []dto.FriendRequest
	for _, request := range r {
		var friendID uuid.UUID
		if userID == request.RecipientID {
			friendID = request.SenderID
		} else {
			friendID = request.RecipientID
		}

		friend, err := conn.Queries.GetUserByID(ctx, friendID)
		if err != nil {
			logger.Warn("Failed to find friend", err)
			return dto.GetFriendRequestsAllResponse{}, err
		}
		requests = append(requests, dto.FriendRequest{FromUsername: friend.Username, Date: request.CreatedAt.Time})
	}

	resp := dto.GetFriendRequestsAllResponse{Requests: requests}
	logger.Debug("Got all friend requests")
	return resp, nil
}
