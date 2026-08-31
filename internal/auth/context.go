package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/hardal7/chrono/internal/util/requestctx"
)

func UserID(ctx context.Context) uuid.UUID {
	var id uuid.UUID

	id, ok := ctx.Value(requestctx.UserID).(uuid.UUID)
	if !ok {
		logger.Warn("Failed to fetch userID")
		return uuid.Nil
	}
	return id
}

func AsUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, requestctx.UserID, userID)
}

func SessionID(ctx context.Context) uuid.UUID {
	var id uuid.UUID

	id, ok := ctx.Value(requestctx.SessionID).(uuid.UUID)
	if !ok {
		logger.Warn("Failed to fetch sessionID")
		return uuid.Nil
	}
	return id
}
