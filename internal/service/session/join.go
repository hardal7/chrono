package session

import (
	"context"
	"errors"
	"fmt"

	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/jackc/pgx/v5"
)

func Join(ctx context.Context, r dto.JoinSessionRequest) error {
	userID := middleware.UserID(ctx)

	s, err := db.Queries.GetSessionByNameAndOwnerName(ctx, query.GetSessionByNameAndOwnerNameParams{
		Name:     r.Name,
		Username: r.OwnerUsername,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("Session not found: %w", db.ErrNotFound)
	}

	if err != nil {
		return fmt.Errorf("Failed to find session: %w: %w", db.ErrRunQuery, err)
	}

	if r.Password != s.Password.String {
		return fmt.Errorf("Wrong password for session")
	}

	p, err := db.Queries.GetSessionParticipantsAsUsers(ctx, s.ID)
	if err != nil {
		return fmt.Errorf("Failed to check if session is full: %w: %w", db.ErrRunQuery, err)
	}

	if s.MaxParticipants.Valid && int(s.MaxParticipants.Int32) == len(p) {
		return fmt.Errorf("Session is full")
	}

	if !s.IsActive {
		return fmt.Errorf("Session has expired")
	}

	err = db.Queries.JoinSession(ctx, query.JoinSessionParams{
		UserID:    userID,
		SessionID: s.ID,
	})
	if err != nil {
		return fmt.Errorf("Failed to join session: %w: %w", db.ErrRunQuery, err)
	}

	err = db.Queries.CreateTopic(ctx, query.CreateTopicParams{
		OwnerID: userID,
		Name:    s.Topic.String,
	})
	if err != nil {
		return fmt.Errorf("Failed to create topic of session: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
