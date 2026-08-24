package location

import (
	"context"

	"github.com/google/uuid"
	conn "github.com/hardal7/chrono/internal/db"
	db "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
	"github.com/hardal7/chrono/internal/util/logger"
)

func EditLocationPrivacy(ctx context.Context, r dto.EditLocationPrivacyRequest) error {
	logger.Debug("Updating user location privacy", "hide", r.Hide)
	err := conn.Queries.SetLocationPrivacy(ctx, db.SetLocationPrivacyParams{ID: ctx.Value(middleware.UserID).(uuid.UUID), HideCountry: r.Hide})
	if err != nil {
		logger.Warn("Failed to update user location privacy")
		return err
	}

	logger.Debug("Updated user location privacy", "hide", r.Hide)
	return nil
}
