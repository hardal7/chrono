package user

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	db "github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/jackc/pgx/v5"
)

const (
	friendStatusNone     = "none"
	friendStatusPending  = "pending"
	friendStatusAccepted = "accepted"

	privateCountry = "Private"
)

func GetProfile(ctx context.Context, username string) (dto.GetUserProfileResponse, error) {
	resp := dto.GetUserProfileResponse{}

	user, err := db.Queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return resp, fmt.Errorf("User not found: %w", db.ErrNotFound)
		}
		return resp, fmt.Errorf("Failed to check if user exists: %w: %w", db.ErrRunQuery, err)
	}

	topics, err := db.Queries.GetTopicsByOwner(ctx, user.ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return resp, fmt.Errorf("Failed to get user topics: %w: %w", db.ErrRunQuery, err)
	}
	var bestTopic string
	var mostTime int32
	for _, topic := range topics {
		if topic.TotalTimeTrackedSeconds > mostTime {
			mostTime = topic.TotalTimeTrackedSeconds
			bestTopic = topic.Name
		}
	}

	possibleFriends, err := db.Queries.GetPossibleFriends(ctx, ctx.Value(middleware.UserID).(uuid.UUID))
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return resp, fmt.Errorf("Failed to get user friends: %w: %w", db.ErrRunQuery, err)
	}
	friendStatus := friendStatusNone
	for _, friend := range possibleFriends {
		if user.Username == friend.Username {
			if friend.IsAccepted {
				friendStatus = friendStatusAccepted
			} else {
				friendStatus = friendStatusPending
			}
		}
	}

	var country string
	if user.HideCountry {
		country = privateCountry
	} else {
		country = user.Country.String
	}

	resp = dto.GetUserProfileResponse{
		Username:         user.Username,
		AvatarPath:       filepath.Join(config.AvatarEndpoint, user.ID.String()),
		TotalTimeSeconds: int(user.TotalTimeTrackedSeconds),
		TodayTimeSeconds: int(user.TodayTimeTrackedSeconds),
		Country:          country,
		BestTopic:        bestTopic,
		FriendStatus:     friendStatus,
	}

	return resp, nil
}
