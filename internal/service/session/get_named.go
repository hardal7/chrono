package session

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/hardal7/chrono/internal/auth"
	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/util/apierror"
	"github.com/hardal7/chrono/internal/util/config"
)

func GetNamed(ctx context.Context, r dto.GetSessionNamedRequest) (dto.GetSessionNamedResponse, error) {
	userID := auth.UserID(ctx)
	resp := dto.GetSessionNamedResponse{}

	s, err := db.Queries.GetSessionByNameAndOwnerName(ctx, query.GetSessionByNameAndOwnerNameParams{
		Name:          r.Name,
		OwnerUsername: r.OwnerUsername,
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
		if participant.ID == userID {
			isParticipant = true
		}

		lastOnlineAgo := int(time.Since(participant.LastSeenAt).Minutes())

		participants = append(participants, dto.Participant{
			Name:             participant.Username,
			AvatarPath:       filepath.Join(config.AvatarEndpoint, participant.UserID.String()),
			SessionTime:      int(participant.TotalTimeTrackedSeconds),
			SessionTimeToday: int(participant.TodayTimeTrackedSeconds),
			LastOnlineAgo:    lastOnlineAgo,
		})
	}

	if !isParticipant {
		return resp, fmt.Errorf("Unauthorized session details requested: %w", apierror.ErrUnauthorized)
	}

	expiresAt := &s.ExpiresAt.Time
	if !s.ExpiresAt.Valid {
		expiresAt = nil
	}

	resp = dto.GetSessionNamedResponse{
		Name:              s.Name,
		OwnerUsername:     r.OwnerUsername,
		ExpiresAt:         expiresAt,
		TotalTime:         int(s.TotalTimeTrackedSeconds),
		MaxParticipants:   int(s.MaxParticipants.Int32),
		TotalParticipants: len(participants),
		Participants:      participants,
	}

	return resp, nil
}
