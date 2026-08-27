package session

import (
	"context"
	"fmt"
	"path/filepath"

	db "github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/config"
)

func GetAll(ctx context.Context) (dto.GetSessionsAllResponse, error) {
	userID := middleware.UserID(ctx)
	resp := dto.GetSessionsAllResponse{}

	s, err := db.Queries.GetSessionsAllByFriends(ctx, userID)
	if err != nil {
		return resp, fmt.Errorf("Failed to get all sessions of friends: %w: %w", db.ErrRunQuery, err)
	}

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
		}

		sessions = append(sessions, dto.SessionSelection{
			Name:              session.Name,
			OwnerUsername:     session.OwnerUsername,
			OwnerAvatarPath:   filepath.Join(config.AvatarEndpoint, session.OwnerID.String()),
			TotalTime:         int(session.TotalTimeSeconds),
			ExpiresAt:         session.ExpiresAt.Time,
			MaxParticipants:   int(session.MaxParticipants.Int32),
			TotalParticipants: len(minParticipants),
			Participants:      minParticipants,
		})
	}

	return dto.GetSessionsAllResponse{Sessions: sessions}, nil
}
