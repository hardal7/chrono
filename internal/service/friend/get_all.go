package friend

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	db "github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
)

func GetAll(ctx context.Context) (dto.GetFriendRequestsAllResponse, error) {
	userID := middleware.UserID(ctx)
	resp := dto.GetFriendRequestsAllResponse{}

	r, err := db.Queries.GetFriendRequests(ctx, userID)
	if err != nil {
		return resp, fmt.Errorf("Failed to get all friend requests: %w: %w", db.ErrRunQuery, err)
	}

	var requests []dto.FriendRequest
	for _, request := range r {
		var friendID uuid.UUID
		if userID == request.RecipientID {
			friendID = request.SenderID
		} else {
			friendID = request.RecipientID
		}

		friend, err := db.Queries.GetUserByID(ctx, friendID)
		if err != nil {
			return resp, fmt.Errorf("Failed to find friend %q: %w: %w", friendID, db.ErrRunQuery, err)
		}
		requests = append(requests, dto.FriendRequest{FromUsername: friend.Username, Date: request.CreatedAt.Time})
	}
	resp = dto.GetFriendRequestsAllResponse{Requests: requests}

	return resp, nil
}
