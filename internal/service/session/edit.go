package session

import (
	"context"
	"fmt"

	"github.com/hardal7/chrono/internal/auth"
	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/jackc/pgx/v5/pgtype"
)

func Edit(ctx context.Context, r dto.EditSessionRequest) error {
	userID := auth.UserID(ctx)

	err := db.Queries.UpdateSession(ctx, query.UpdateSessionParams{
		OwnerID:         userID,
		Name:            r.Name,
		NewName:         r.NewName,
		MaxParticipants: pgtype.Int4{Int32: int32(r.NewMaxParticipants), Valid: r.NewMaxParticipants != 0},
		ExpiresAt:       pgtype.Timestamptz{Time: r.NewExpiresAt, Valid: !r.NewExpiresAt.IsZero()},
	})
	if err != nil {
		return fmt.Errorf("Failed to update session: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
