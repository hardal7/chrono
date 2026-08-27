package session

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
)

func Kick(ctx context.Context, r dto.KickFromSessionRequest) error {
	err := db.Queries.KickFromSession(ctx, query.KickFromSessionParams{
		OwnerID:             ctx.Value(middleware.UserID).(uuid.UUID),
		Name:                r.SessionName,
		ParticipantUsername: r.ParticipantUsername,
	})
	if err != nil {
		return fmt.Errorf("Failed to kick user from session: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
