package session

import (
	"context"
	"fmt"

	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/util/requestctx"
)

func Leave(ctx context.Context, r dto.LeaveSessionRequest) error {
	userID := requestctx.GetUserID(ctx)

	err := db.Queries.LeaveSession(ctx, query.LeaveSessionParams{UserID: userID, Name: r.Name, OwnerUsername: r.OwnerUsername})
	if err != nil {
		return fmt.Errorf("leave session: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
