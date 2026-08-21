package user

import (
	"context"
	"path/filepath"

	"github.com/google/uuid"
	conn "github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5"
)

func GetProfile(ctx context.Context, username string) (dto.GetUserProfileResponse, error) {
	logger.Debug("Getting user profile", "username", username)
	user, err := conn.Queries.GetUserByUsername(ctx, username)
	if err != nil {
		if err == pgx.ErrNoRows {
			logger.Debug("User not found")
			return dto.GetUserProfileResponse{}, err
		} else {
			logger.Warn("Failed to check if user exists", err)
			return dto.GetUserProfileResponse{}, err
		}
	}

	topics, err := conn.Queries.GetTopicsByOwner(ctx, user.ID)
	if err != nil && err != pgx.ErrNoRows {
		logger.Warn("Failed to get user topics", err)
		return dto.GetUserProfileResponse{}, err
	}
	var bestTopic string
	var mostTime int32
	for _, topic := range topics {
		if topic.TotalTimeTrackedSeconds > mostTime {
			mostTime = topic.TotalTimeTrackedSeconds
			bestTopic = topic.Name
		}
	}

	possibleFriends, err := conn.Queries.GetPossibleFriends(ctx, ctx.Value(middleware.UserID).(uuid.UUID))
	if err != nil && err != pgx.ErrNoRows {
		logger.Warn("Failed to get user friends", err)
		return dto.GetUserProfileResponse{}, err
	}
	friendStatus := "none"
	for _, friend := range possibleFriends {
		if user.Username == friend.Username {
			if friend.IsAccepted {
				friendStatus = "accepted"
			} else {
				friendStatus = "pending"
			}
		}
	}

	var country string
	if user.HideCountry {
		country = "Private"
	} else {
		country = user.Country.String
	}

	resp := dto.GetUserProfileResponse{
		Username:         user.Username,
		AvatarPath:       filepath.Join(config.AvatarEndpoint, user.ID.String()),
		TotalTimeSeconds: int(user.TotalTimeTrackedSeconds),
		TodayTimeSeconds: int(user.TodayTimeTrackedSeconds),
		Country:          country,
		BestTopic:        bestTopic,
		FriendStatus:     friendStatus,
	}

	logger.Debug("Got user profile", "username", username)
	return resp, nil
}
