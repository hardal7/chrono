package location

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/middleware"
)

func EditLocationPrivacy(ctx context.Context, r dto.EditLocationPrivacyRequest) error {
	err := db.Queries.SetLocationPrivacy(ctx, query.SetLocationPrivacyParams{ID: ctx.Value(middleware.UserID).(uuid.UUID), HideCountry: r.Hide})
	if err != nil {
		return fmt.Errorf("Failed to update user location privacy: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
