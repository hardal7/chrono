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
)

func GetNamed(ctx context.Context, r dto.GetSessionNamedRequest) (dto.GetSessionNamedResponse, error) {
	logger.Debug("Getting session", "sessionName", r.Name)
	s, err := conn.Queries.GetSessionByNameAndOwnerName(ctx, db.GetSessionByNameAndOwnerNameParams{
		Name:     r.Name,
		Username: r.OwnerUsername,
	})
	if err != nil {
		logger.Debug("Failed to get session by username", err)
		return dto.GetSessionNamedResponse{}, err
	}

	p, err := conn.Queries.GetSessionParticipants(ctx, s.ID)
	if err != nil {
		logger.Warn("Failed to get participants of the session", "sessionName", s.Name, "error", err)
		return dto.GetSessionNamedResponse{}, err
	}
	participants := []dto.Participant{}
	isParticipant := false
	for _, participant := range p {
		u, err := conn.Queries.GetUserByID(ctx, participant.ID)
		if err != nil {
			logger.Warn("Failed to get user participant", err)
			return dto.GetSessionNamedResponse{}, err
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
		logger.Debug("Cannot get details for unauthorized session")
		return dto.GetSessionNamedResponse{}, errors.New("unauthorized session details requested")
	}

	resp := dto.GetSessionNamedResponse{
		Name:                s.Name,
		OwnerUsername:       r.OwnerUsername,
		ExpiresAt:           s.ExpiresAt.Time,
		MaxParticipants:     int(s.MaxParticipants.Int32),
		TotalParticipants:   len(participants),
		CurrentParticipants: participants,
	}
	logger.Debug("Got session", "sessionName", r.Name)
	return resp, nil
}
