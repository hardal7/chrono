package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/util/requestctx"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	scopeFriends = "friends"
	scopeLocal   = "local"
	scopeGlobal  = "global"

	TopUsersKey = "top_users"
)

func GetTopUsers(ctx context.Context, r dto.GetTopUsersRequest) (dto.GetTopUsersResponse, error) {
	userID := requestctx.GetUserID(ctx)
	var users []query.User
	var err error
	resp := dto.GetTopUsersResponse{}
	matchName := pgtype.Text{String: r.MatchName, Valid: true}

	switch r.Scope {
	case scopeFriends:
		users, err = db.Queries.GetTopFriends(ctx, query.GetTopFriendsParams{
			ID:        userID,
			Cursor:    int32(r.Cursor),
			Limit:     int32(r.Limit),
			MatchName: matchName,
		})
		if err != nil {
			break
		}

		var user query.User
		user, err = db.Queries.GetUserByID(ctx, userID)
		if err != nil {
			return resp, fmt.Errorf("retrieve user: %w: %w", db.ErrRunQuery, err)
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
		var data []byte
		data, err = db.RDB.Get(ctx, TopUsersKey).Bytes()
		if err != nil {
			break
		}

		err = json.Unmarshal(data, &resp)
		if err != nil {
			break
		}

	default:
		return resp, errors.New("invalid scope queried")
	}
	if err != nil {
		return resp, fmt.Errorf("get users: %w: %w", db.ErrRunQuery, err)
	}

	// TODO: Rank changes

	return resp, nil
}
