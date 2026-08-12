package session

import (
	"context"
	"time"

	"github.com/google/uuid"
	conn "github.com/hardal7/chrono/internal/db"
	db "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
)

func Join(ctx context.Context, r dto.JoinSessionRequest) error {
	logger.Info("Joining session", "sessionName", r.Name)

	s, err := conn.Queries.GetSessionByNameAndOwnerName(ctx, db.GetSessionByNameAndOwnerNameParams{
		Name:     r.Name,
		Username: r.OwnerUsername,
	})
	if err != nil {
		logger.Error("Failed to find session", err)
		return err
	}
	if r.Password != s.Password.String {
		logger.Error("Wrong password for session", err)
		return err
	}
	p, err := conn.Queries.GetSessionParticipantsBySessionID(ctx, s.ID)
	if err != nil {
		logger.Error("Failed to check if session is full", err)
		return err
	}
	if int(s.MaxParticipants.Int32) == len(p) {
		logger.Error("Session is full")
		return err
	}
	if !s.IsActive {
		logger.Error("Session has expired")
		return err
	}

	err = conn.Queries.JoinSession(ctx, db.JoinSessionParams{
		UserID:     ctx.Value(middleware.UserID).(uuid.UUID),
		SessionID:  s.ID,
		LastSeenAt: time.Now(),
	})
	if err != nil {
		logger.Error("Failed to join session", err)
		return err
	}
	logger.Info("Joined session", "sessionName", r.Name)
	return nil
}
