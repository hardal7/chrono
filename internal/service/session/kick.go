package session

import (
	"context"
	"fmt"

	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/util/requestctx"
)

func Kick(ctx context.Context, r dto.KickFromSessionRequest) error {
	userID := requestctx.GetUserID(ctx)

	err := db.Queries.KickFromSession(ctx, query.KickFromSessionParams{
		OwnerID:             userID,
		Name:                r.SessionName,
		ParticipantUsername: r.ParticipantUsername,
	})
	if err != nil {
		return fmt.Errorf("kick user from session: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
