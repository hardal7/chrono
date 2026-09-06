package location

import (
	"context"
	"fmt"

	"github.com/hardal7/chrono/internal/auth"
	db "github.com/hardal7/chrono/internal/db"
	query "github.com/hardal7/chrono/internal/db/sqlc"
	"github.com/hardal7/chrono/internal/dto"
)

func EditLocationPrivacy(ctx context.Context, r dto.EditLocationPrivacyRequest) error {
	userID := auth.UserID(ctx)

	err := db.Queries.SetLocationPrivacy(ctx, query.SetLocationPrivacyParams{ID: userID, HideCountry: r.Hide})
	if err != nil {
		return fmt.Errorf("Failed to update user location privacy: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
