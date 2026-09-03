package user

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/hardal7/chrono/internal/auth"
	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	scopeFriends = "friends"
	scopeLocal   = "local"
	scopeGlobal  = "global"
)

func GetTopUsers(ctx context.Context, r dto.GetTopUsersRequest) (dto.GetTopUsersResponse, error) {
	userID := auth.UserID(ctx)
	var users []query.User
	var err error
	resp := dto.GetTopUsersResponse{}
	matchName := pgtype.Text{String: r.MatchName, Valid: r.MatchName != ""}

	switch r.Scope {
	case scopeFriends:
		users, err = db.Queries.GetTopFriends(ctx, query.GetTopFriendsParams{
			ID:        userID,
			Cursor:    int32(r.Cursor),
			Limit:     int32(r.Limit),
			MatchName: matchName,
		})
		user, err := db.Queries.GetUserByID(ctx, userID)
		if err != nil {
			return resp, fmt.Errorf("Failed to retrieve user: %w: %w", db.ErrRunQuery, err)
		}
		users = append(users, user)

	case scopeLocal:
		users, err = db.Queries.GetTopUsersLocal(ctx, query.GetTopUsersLocalParams{
			ID:        userID,
			Cursor:    int32(r.Cursor),
			Limit:     int32(r.Limit),
			MatchName: matchName,
		})

	case scopeGlobal:
		// TODO: Cache this with redis (update on 1m?)
		users, err = db.Queries.GetTopUsers(ctx, query.GetTopUsersParams{
			Cursor:    int32(r.Cursor),
			Limit:     int32(r.Limit),
			MatchName: matchName,
		})

	default:
		return resp, fmt.Errorf("Invalid scope queried")
	}
	if err != nil {
		return resp, fmt.Errorf("Failed to get users: %w: %w", db.ErrRunQuery, err)
	}

	// TODO: Rank changes
	for i, user := range users {
		resp.Users = append(resp.Users, dto.TopUser{
			Rank:       i + 1,
			RankChange: 7,
			Username:   user.Username,
			TotalTime:  int(user.TotalTimeTrackedSeconds),
			TodayTime:  int(user.TodayTimeTrackedSeconds),
			AvatarPath: filepath.Join(config.AvatarEndpoint, user.ID.String()),
		})
	}

	return resp, nil
}
