package session

import (
	"context"

	"github.com/google/uuid"
	conn "github.com/hardal7/chrono/internal/db"
	db "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5/pgtype"
)

func Edit(ctx context.Context, r dto.EditSessionRequest) error {
	s, err := conn.Queries.GetSessionByNameAndOwnerID(ctx, db.GetSessionByNameAndOwnerIDParams{
		Name:    r.Name,
		OwnerID: ctx.Value(middleware.UserID).(uuid.UUID),
	})
	if err != nil {
		return err
	}

	logger.Info("Editing session", "sessionName", s.Name)
	if r.Delete {
		logger.Info("Deleting session", "sessionName", s.Name)
		err := conn.Queries.DeleteSession(ctx, s.ID)
		if err != nil {
			logger.Error("Failed to delete session", err)
			return err
		}
		logger.Info("Deleted session")
		return nil
	}

	if r.NewName == "" {
		s.Name = r.NewName
	}
	err = conn.Queries.UpdateSession(ctx, db.UpdateSessionParams{
		ID:              s.ID,
		Name:            s.Name,
		MaxParticipants: pgtype.Int4{Int32: int32(r.NewMaxParticipants), Valid: r.NewMaxParticipants != 0},
		Password:        pgtype.Text{String: r.NewPassword, Valid: r.NewPassword != ""},
		ExpiresAt:       pgtype.Timestamptz{Time: r.NewExpiresAt, Valid: !r.NewExpiresAt.IsZero()},
	})
	if err != nil {
		logger.Error("Failed to update session", err)
		return err
	}
	logger.Info("Changed session details")
	return nil
}
