package session

import (
	"context"
	"fmt"

	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/util/requestctx"
)

func Delete(ctx context.Context, r dto.DeleteSessionRequest) error {
	userID := requestctx.GetUserID(ctx)

	err := db.Queries.DeleteSession(ctx, query.DeleteSessionParams{
		OwnerID: userID,
		Name:    r.Name,
	})
	if err != nil {
		return fmt.Errorf("delete session: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
