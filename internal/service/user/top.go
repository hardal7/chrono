package user

import (
	"context"

	conn "github.com/hardal7/chrono/internal/db"
	db "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/util/logger"
)

func GetTopUsers(ctx context.Context, r dto.GetTopUsersRequest) (dto.GetTopUsersResponse, error) {
	users, err := conn.Queries.GetTopUsers(ctx, db.GetTopUsersParams{
		TimeTrackedSeconds: int32(r.Cursor),
		Limit:              int32(r.Limit),
	})
	if err != nil {
		logger.Error("Failed to get users", err)
		return dto.GetTopUsersResponse{}, err
	} else {
		resp := dto.GetTopUsersResponse{}
		for i, v := range users {
			resp.Users[i] = dto.TopUser{
				Rank:      i + 1,
				Username:  v.Username,
				TotalTime: int(v.TimeTrackedSeconds),
				// TODO
				// TodayTime: int(v.TodayTrackedSeconds),
			}
		}
		return resp, nil
	}
}
