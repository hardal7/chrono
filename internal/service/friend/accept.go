package friend

import (
	"context"
	"fmt"

	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/util/requestctx"
)

func AcceptRequest(ctx context.Context, r dto.AcceptFriendRequestRequest) error {
	userID := requestctx.GetUserID(ctx)

	err := db.Queries.AcceptFriendRequest(ctx, query.AcceptFriendRequestParams{
		RecipientID: userID,
		Username:    r.Username,
	})
	if err != nil {
		return fmt.Errorf("accept friend request: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
