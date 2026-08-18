package session

import (
	"context"
)

// TODO: Check if session has expired before returning,
// If expired, issue a deletion query of expired sessions

func GetAll(ctx context.Context) error {
	return nil
}
