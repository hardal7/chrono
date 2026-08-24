package session

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5/pgtype"
)

func Create(ctx context.Context, r dto.CreateSessionRequest) error {
	err := db.Queries.CreateSession(ctx, query.CreateSessionParams{
		Name:            r.Name,
		OwnerID:         ctx.Value(middleware.UserID).(uuid.UUID),
		MaxParticipants: pgtype.Int4{Int32: int32(r.MaxParticipants), Valid: r.MaxParticipants != 0},
		Password:        pgtype.Text{String: r.Password, Valid: r.Password != ""},
		ExpiresAt:       pgtype.Timestamptz{Time: r.ExpiresAt, Valid: !r.ExpiresAt.IsZero()},
		Topic:           pgtype.Text{String: r.Topic, Valid: r.Topic != ""},
	})
	if err != nil {
		return fmt.Errorf("Failed to create session: %w: %w", db.ErrRunQuery, err)
	}

	u, _ := db.Queries.GetUserByID(ctx, ctx.Value(middleware.UserID).(uuid.UUID))
	err = Join(ctx, dto.JoinSessionRequest{Name: r.Name, Password: r.Password, OwnerUsername: u.Username})
	if err != nil {
		logger.Warn("Failed to join own session")
	}

	return nil
}
