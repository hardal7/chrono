package session

import (
	"context"
	"fmt"

	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/jackc/pgx/v5/pgtype"
)

func Edit(ctx context.Context, r dto.EditSessionRequest) error {
	userID := middleware.UserID(ctx)

	s, err := db.Queries.GetSessionByNameAndOwnerID(ctx, query.GetSessionByNameAndOwnerIDParams{
		Name:    r.Name,
		OwnerID: userID,
	})
	if err != nil {
		return fmt.Errorf("Failed to get session: %w: %w", db.ErrRunQuery, err)
	}

	if r.NewName == "" {
		s.Name = r.NewName
	}
	err = db.Queries.UpdateSession(ctx, query.UpdateSessionParams{
		ID:              s.ID,
		Name:            s.Name,
		MaxParticipants: pgtype.Int4{Int32: int32(r.NewMaxParticipants), Valid: r.NewMaxParticipants != 0},
		Password:        pgtype.Text{String: r.NewPassword, Valid: r.NewPassword != ""},
		ExpiresAt:       pgtype.Timestamptz{Time: r.NewExpiresAt, Valid: !r.NewExpiresAt.IsZero()},
	})
	if err != nil {
		return fmt.Errorf("Failed to update session: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
