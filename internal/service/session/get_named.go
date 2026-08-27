package session

import (
	"context"
	"fmt"

	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/apierror"
)

func GetNamed(ctx context.Context, r dto.GetSessionNamedRequest) (dto.GetSessionNamedResponse, error) {
	userID := middleware.UserID(ctx)
	resp := dto.GetSessionNamedResponse{}

	s, err := db.Queries.GetSessionByNameAndOwnerName(ctx, query.GetSessionByNameAndOwnerNameParams{
		Name:     r.Name,
		Username: r.OwnerUsername,
	})
	if err != nil {
		return resp, fmt.Errorf("Failed to get session by username: %w: %w", db.ErrRunQuery, err)
	}

	p, err := db.Queries.GetSessionParticipantsAsUsers(ctx, s.ID)
	if err != nil {
		return resp, fmt.Errorf("Failed to get participants of the session: %w: %w", db.ErrRunQuery, err)
	}

	participants := []dto.Participant{}
	isParticipant := false
	for _, participant := range p {
		if participant.ID == userID {
			isParticipant = true
		}

		participants = append(participants, dto.Participant{
			Name:             participant.Username,
			AvatarPath:       participant.ID.String(),
			SessionTime:      int(participant.TotalTimeTrackedSeconds),
			SessionTimeToday: int(participant.TodayTimeTrackedSeconds),
			LastOnline:       int(participant.LastSeenAt.Unix()),
		})
	}

	if !isParticipant {
		return resp, fmt.Errorf("Unauthorized session details requested: %w", apierror.ErrUnauthorized)
	}

	resp = dto.GetSessionNamedResponse{
		Name:                s.Name,
		OwnerUsername:       r.OwnerUsername,
		ExpiresAt:           s.ExpiresAt.Time,
		MaxParticipants:     int(s.MaxParticipants.Int32),
		TotalParticipants:   len(participants),
		CurrentParticipants: participants,
	}

	return resp, nil
}
