package requestctx

import (
	"context"
	"uuid"

	"github.com/hardal7/chrono/internal/util/logger"
)

func GetUserID(ctx context.Context) uuid.UUID {
	var id uuid.UUID

	id, ok := ctx.Value(UserID).(uuid.UUID)
	if !ok {
		logger.Warn("fetch userID")
		return uuid.Nil()
	}
	return id
}

func AsUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, UserID, userID)
}

func GetSessionID(ctx context.Context) string {
	var sessionID string

	sessionID, ok := ctx.Value(SessionID).(string)
	if !ok {
		logger.Warn("fetch sessionID")
		return ""
	}
	return sessionID
}

func GetRequestID(ctx context.Context) string {
	var requestID string

	requestID, ok := ctx.Value(RequestID).(string)
	if !ok {
		logger.Warn("fetch requestID")
		return ""
	}
	return requestID
}
