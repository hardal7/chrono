package user

import (
	"context"
	"fmt"

	"github.com/hardal7/chrono/internal/auth"
	"github.com/hardal7/chrono/internal/db"
)

func Logout(ctx context.Context) error {
	sessionID := auth.SessionID(ctx)

	err := db.Queries.DeleteSessionToken(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("Failed to delete session token: %w: %w", db.ErrRunQuery, err)
	}

	return nil
}
