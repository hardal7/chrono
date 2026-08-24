package session

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
)

func GetNamed(ctx context.Context, r dto.GetSessionNamedRequest) (dto.GetSessionNamedResponse, error) {
	resp := dto.GetSessionNamedResponse{}

	s, err := db.Queries.GetSessionByNameAndOwnerName(ctx, query.GetSessionByNameAndOwnerNameParams{
		Name:     r.Name,
		Username: r.OwnerUsername,
	})
	if err != nil {
		return resp, fmt.Errorf("Failed to get session by username: %w: %w", db.ErrRunQuery, err)
	}

	p, err := db.Queries.GetSessionParticipants(ctx, s.ID)
	if err != nil {
		return resp, fmt.Errorf("Failed to get participants of the session: %w: %w", db.ErrRunQuery, err)
	}

	participants := []dto.Participant{}
	isParticipant := false
	for _, participant := range p {
		u, err := db.Queries.GetUserByID(ctx, participant.ID)
		if err != nil {
			return resp, fmt.Errorf("Failed to get user participant: %w: %w", db.ErrRunQuery, err)
		}
		if u.ID == ctx.Value(middleware.UserID).(uuid.UUID) {
			isParticipant = true
		}

		participants = append(participants, dto.Participant{
			Name:             u.Username,
			AvatarPath:       participant.ID.String(),
			SessionTime:      int(participant.TotalTimeTrackedSeconds),
			SessionTimeToday: int(participant.TodayTimeTrackedSeconds),
			LastOnline:       int(participant.LastSeenAt.Unix()),
		})
	}

	if !isParticipant {
		return resp, fmt.Errorf("Unauthorized session details requested")
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
