package auth

import (
	"context"

	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/hardal7/chrono/internal/util/requestctx"
	"uuid"
)

func UserID(ctx context.Context) uuid.UUID {
	var id uuid.UUID

	id, ok := ctx.Value(requestctx.UserID).(uuid.UUID)
	if !ok {
		logger.Warn("Failed to fetch userID")
		return uuid.Nil()
	}
	return id
}

func AsUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, requestctx.UserID, userID)
}

func Session(ctx context.Context) string {
	var session string

	session, ok := ctx.Value(requestctx.SessionID).(string)
	if !ok {
		logger.Warn("Failed to fetch session")
		return ""
	}
	return session
}
