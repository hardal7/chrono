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
	users, err := conn.Queries.GetTopUsers(ctx, db.GetTopUsersParams{
		TotalTimeTrackedSeconds: int32(r.Cursor),
		Limit:                   int32(r.Limit),
	})
	if err != nil {
		logger.Error("Failed to get users", err)
		return dto.GetTopUsersResponse{}, err
	} else {
		resp := dto.GetTopUsersResponse{}
		for i, v := range users {
			avatar, err := conn.Queries.GetAvatarFromUserID(ctx, ctx.Value(middleware.UserID).(uuid.UUID))
			avatarPath := avatar.ID.String()
			if err != nil {
				logger.Warn("Failed to get avatar of user", "userRank", i+1)
				avatarPath = "empty"
			}
			logger.Warn("user", v.Username, "time", v.TotalTimeTrackedSeconds)
			resp.Users = append(resp.Users, dto.TopUser{
				Rank:       i + 1,
				Username:   v.Username,
				TotalTime:  int(v.TotalTimeTrackedSeconds),
				TodayTime:  int(v.TodayTimeTrackedSeconds),
				AvatarPath: "/avatar/" + avatarPath,
			})
		}
		return resp, nil
	}
}
