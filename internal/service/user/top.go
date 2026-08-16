package user

import (
	"context"

	"github.com/google/uuid"
	conn "github.com/hardal7/chrono/internal/db"
	db "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
)

func GetTopUsers(ctx context.Context, r dto.GetTopUsersRequest) (dto.GetTopUsersResponse, error) {
	logger.Info("Getting top users", "friends_only", r.FriendsOnly)
	var users []db.User
	var err error
	if r.FriendsOnly {
		users, err = conn.Queries.GetTopFriends(ctx, db.GetTopFriendsParams{
			TotalTimeTrackedSeconds: int32(r.Cursor),
			Limit:                   int32(r.Limit),
		})
	} else {
		users, err = conn.Queries.GetTopUsers(ctx, db.GetTopUsersParams{
			TotalTimeTrackedSeconds: int32(r.Cursor),
			Limit:                   int32(r.Limit),
		})
	}
	if err != nil {
		logger.Error("Failed to get users", err)
		return dto.GetTopUsersResponse{}, err
	}

	resp := dto.GetTopUsersResponse{}
	for i, user := range users {
		avatarPath := ctx.Value(middleware.UserID).(uuid.UUID).String()
		logger.Warn("user", user.Username, "time", user.TotalTimeTrackedSeconds)
		resp.Users = append(resp.Users, dto.TopUser{
			Rank:       i + 1,
			Username:   user.Username,
			TotalTime:  int(user.TotalTimeTrackedSeconds),
			TodayTime:  int(user.TodayTimeTrackedSeconds),
			AvatarPath: "/avatar/" + avatarPath,
		})
	}
	logger.Info("Got top users")
	return resp, nil
}
