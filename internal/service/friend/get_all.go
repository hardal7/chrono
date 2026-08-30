package friend

import (
	"context"
	"fmt"

	db "github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/auth"
)

func GetAll(ctx context.Context) (dto.GetFriendRequestsAllResponse, error) {
	userID := auth.UserID(ctx)
	resp := dto.GetFriendRequestsAllResponse{}

	reqs, err := db.Queries.GetFriendRequests(ctx, userID)
	if err != nil {
		return resp, fmt.Errorf("Failed to get all friend requests: %w: %w", db.ErrRunQuery, err)
	}

	var requests []dto.FriendRequest
	for _, request := range reqs {
		requests = append(requests, dto.FriendRequest{FromUsername: request.Username, Date: request.CreatedAt.Time})
	}
	resp = dto.GetFriendRequestsAllResponse{Requests: requests}

	return resp, nil
}
