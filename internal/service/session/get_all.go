package session

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	db "github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
)

func GetAll(ctx context.Context) (dto.GetSessionsAllResponse, error) {
	resp := dto.GetSessionsAllResponse{}

	s, err := db.Queries.GetSessionsAllByFriends(ctx, ctx.Value(middleware.UserID).(uuid.UUID))
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
			minParticipants = append(minParticipants, dto.MinParticipant{Name: participant.Username, AvatarPath: participant.ID.String()})
		}

		sessions = append(sessions, dto.SessionSelection{
			Name:              session.Name,
			MaxParticipants:   int(session.MaxParticipants.Int32),
			OwnerUsername:     session.OwnerUsername,
			TotalParticipants: len(minParticipants),
			Participants:      minParticipants,
		})
	}

	return dto.GetSessionsAllResponse{Sessions: sessions}, nil
}
