package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	conn "github.com/hardal7/chrono/internal/db"
	db "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetTopUsers(ctx context.Context, r dto.GetTopUsersRequest) (dto.GetTopUsersResponse, error) {
	logger.Info("Getting top users", "scope", r.Scope)
	var users []db.User
	var err error
	matchName := pgtype.Text{String: r.MatchName, Valid: true}

	switch r.Scope {
	case "friends":
		users, err = conn.Queries.GetTopFriends(ctx, db.GetTopFriendsParams{
			ID:                      ctx.Value(middleware.UserID).(uuid.UUID),
			TotalTimeTrackedSeconds: int32(r.Cursor),
			Limit:                   int32(r.Limit),
			MatchName:               matchName,
		})
		user, err := conn.Queries.GetUserByID(ctx, ctx.Value(middleware.UserID).(uuid.UUID))
		if err != nil {
			logger.Error("Failed to retrieve user", err)
			return dto.GetTopUsersResponse{}, err
		}
		users = append(users, user)
	case "local":
		users, err = conn.Queries.GetTopUsersLocal(ctx, db.GetTopUsersLocalParams{
			ID:                      ctx.Value(middleware.UserID).(uuid.UUID),
			TotalTimeTrackedSeconds: int32(r.Cursor),
			Limit:                   int32(r.Limit),
			MatchName:               matchName,
		})
	case "global":
		users, err = conn.Queries.GetTopUsers(ctx, db.GetTopUsersParams{
			TotalTimeTrackedSeconds: int32(r.Cursor),
			Limit:                   int32(r.Limit),
			MatchName:               matchName,
		})
	default:
		logger.Error("Invalid scope queried", "scope", r.Scope, err)
		return dto.GetTopUsersResponse{}, errors.New("invalid scope")
	}
	if err != nil {
		logger.Error("Failed to get users", err)
		return dto.GetTopUsersResponse{}, err
	}

	resp := dto.GetTopUsersResponse{}
	for i, user := range users {
		avatarPath := ctx.Value(middleware.UserID).(uuid.UUID).String()
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
