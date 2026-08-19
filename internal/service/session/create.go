package session

import (
	"context"
	"errors"

	"github.com/google/uuid"
	conn "github.com/hardal7/chrono/internal/db"
	db "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func Create(ctx context.Context, r dto.CreateSessionRequest) error {
	logger.Debug("Creating session", "sessionName", r.Name)

	_, err := conn.Queries.GetSessionByNameAndOwnerID(ctx, db.GetSessionByNameAndOwnerIDParams{
		Name:    r.Name,
		OwnerID: ctx.Value(middleware.UserID).(uuid.UUID),
	})
	if err != pgx.ErrNoRows {
		if err == nil {
			logger.Debug("Session with name exists")
			return errors.New("session with name exists")
		} else {
			logger.Warn("Failed to check if session is duplicate", err)
			return err
		}
	}

	err = conn.Queries.CreateSession(ctx, db.CreateSessionParams{
		Name:            r.Name,
		OwnerID:         ctx.Value(middleware.UserID).(uuid.UUID),
		MaxParticipants: pgtype.Int4{Int32: int32(r.MaxParticipants), Valid: r.MaxParticipants != 0},
		Password:        pgtype.Text{String: r.Password, Valid: r.Password != ""},
		ExpiresAt:       pgtype.Timestamptz{Time: r.ExpiresAt, Valid: !r.ExpiresAt.IsZero()},
		Topic:           pgtype.Text{String: r.Topic, Valid: r.Topic != ""},
	})
	if err != nil {
		logger.Debug("Failed to create session", err)
		return err
	}
	logger.Debug("Created session", "sessionName", r.Name)

	u, _ := conn.Queries.GetUserByID(ctx, ctx.Value(middleware.UserID).(uuid.UUID))
	logger.Debug("Joining own session")
	err = Join(ctx, dto.JoinSessionRequest{Name: r.Name, Password: r.Password, OwnerUsername: u.Username})
	if err != nil {
		logger.Warn("Failed to join own session")
	}

	return nil
}
