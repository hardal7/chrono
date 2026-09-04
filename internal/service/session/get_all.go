package session

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/hardal7/chrono/internal/auth"
	db "github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/util/config"
)

func GetAll(ctx context.Context) (dto.GetSessionsAllResponse, error) {
	userID := auth.UserID(ctx)
	resp := dto.GetSessionsAllResponse{}

	s, err := db.Queries.GetSessionsAll(ctx, userID)
	if err != nil {
		return resp, fmt.Errorf("Failed to get all sessions of friends: %w: %w", db.ErrRunQuery, err)
	}

	joined := false

	sessions := []dto.SessionSelection{}
	for _, session := range s {
		p, err := db.Queries.GetSessionParticipantsAsUsers(ctx, session.ID)
		if err != nil {
			return resp, fmt.Errorf("Failed to get participants of session %q: %w: %w", session.Name, db.ErrRunQuery, err)
		}

		minParticipants := []dto.MinParticipant{}
		for _, participant := range p {
			minParticipants = append(minParticipants, dto.MinParticipant{
				Name:       participant.Username,
				AvatarPath: filepath.Join(config.AvatarEndpoint, participant.ID.String()),
			})

			if participant.ID == userID {
				joined = true
			}
		}

		expiresAt := &session.ExpiresAt.Time
		if !session.ExpiresAt.Valid {
			expiresAt = nil
		}

		sessions = append(sessions, dto.SessionSelection{
			Name:              session.Name,
			OwnerUsername:     session.OwnerUsername,
			Joined:            joined,
			OwnerAvatarPath:   filepath.Join(config.AvatarEndpoint, session.OwnerID.String()),
			TotalTime:         int(session.TotalTimeTrackedSeconds),
			ExpiresAt:         expiresAt,
			MaxParticipants:   int(session.MaxParticipants.Int32),
			TotalParticipants: len(minParticipants),
			Participants:      minParticipants,
		})
	}

	return dto.GetSessionsAllResponse{Sessions: sessions}, nil
}
