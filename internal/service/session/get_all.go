package session

import (
	"context"
	"time"

	"github.com/google/uuid"
	conn "github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
)

func GetAll(ctx context.Context) (dto.GetSessionsAllResponse, error) {
	logger.Debug("Getting all sessions of friends")
	s, err := conn.Queries.GetSessionsAllByFriends(ctx, ctx.Value(middleware.UserID).(uuid.UUID))
	if err != nil {
		logger.Debug("Failed to get all sessions of friends")
		return dto.GetSessionsAllResponse{}, err
	}

	sessions := []dto.SessionSelection{}
	for _, session := range s {
		if !session.ExpiresAt.Time.After(time.Now()) {
			logger.Debug("Deleting expired session", "session", session.Name)
			err := conn.Queries.DeleteSession(ctx, session.ID)
			if err != nil {
				logger.Warn("Failed to delete expired session", "sessionName", session.Name, "error", err)
			}
		}

		owner, err := conn.Queries.GetUserByID(ctx, session.OwnerID)
		if err != nil {
			logger.Warn("Failed to get the owner of the session", "sessionName", session.Name, "error", err)
			return dto.GetSessionsAllResponse{}, err
		}

		p, err := conn.Queries.GetSessionParticipants(ctx, session.ID)
		if err != nil {
			logger.Warn("Failed to get participants of the session", "sessionName", session.Name, "error", err)
			return dto.GetSessionsAllResponse{}, err
		}
		minParticipants := []dto.MinParticipant{}
		for _, participant := range p {
			u, err := conn.Queries.GetUserByID(ctx, participant.ID)
			if err != nil {
				logger.Warn("Failed to get participant", "sessionName", session.Name, "error", err, "userID", participant.ID)
				return dto.GetSessionsAllResponse{}, err
			}
			minParticipants = append(minParticipants, dto.MinParticipant{Name: u.Username, AvatarPath: participant.ID.String()})
		}

		sessions = append(sessions, dto.SessionSelection{
			Name:              session.Name,
			MaxParticipants:   int(session.MaxParticipants.Int32),
			OwnerUsername:     owner.Username,
			TotalParticipants: len(minParticipants),
			Participants:      minParticipants,
		})
	}

	logger.Debug("Got all sessions of friends")
	return dto.GetSessionsAllResponse{Sessions: sessions}, nil
}
