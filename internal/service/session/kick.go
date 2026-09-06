package session

import (
	"context"
	"fmt"

	"github.com/hardal7/chrono/internal/auth"
	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
)

func Kick(ctx context.Context, r dto.KickFromSessionRequest) error {
	userID := auth.UserID(ctx)

	err := db.Queries.KickFromSession(ctx, query.KickFromSessionParams{
		OwnerID:             userID,
		Name:                r.SessionName,
		ParticipantUsername: r.ParticipantUsername,
	})
	if err != nil {
		return fmt.Errorf("Failed to kick user from session: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
