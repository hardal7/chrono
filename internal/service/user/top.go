package user

import (
	"context"

	conn "github.com/hardal7/chrono/internal/db"
	db "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
)

func GetTopUsers(ctx context.Context, r dto.GetTopUsersRequest) (dto.GetTopUsersResponse, error) {
	users, err := conn.Queries.GetTopUsers(ctx, db.GetTopUsersParams{
		TimeTrackedSeconds: int32(r.Cursor),
		Limit:              int32(r.Limit),
	})
	if err != nil {
		return dto.GetTopUsersResponse{}, err
	} else {
		resp := dto.GetTopUsersResponse{}
		for i, v := range users {
			resp.Usernames[i] = v.Username
			resp.TotalTimes[i] = int(v.TimeTrackedSeconds)
			// TODO
			// resp.TodayTimes[i] = conn.Queries.GetTimeToday
		}
		return resp, nil
	}
}
