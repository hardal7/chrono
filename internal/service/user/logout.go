package user

import (
	"context"

	"github.com/hardal7/chrono/internal/auth"
)

func Logout(ctx context.Context) error {
	err := auth.DeleteSession(ctx)
	return err
}
