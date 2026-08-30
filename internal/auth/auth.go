package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
)

const (
	AuthHeader = "Authorization"
	Bearer     = "Bearer "

	userID    = "userID"
	sessionID = "sessionID"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Authenticating user")

		token := r.Header.Get(AuthHeader)
		if token == "" {
			logger.Debug("No token provided")
			http.Error(w, "No token provided", http.StatusUnauthorized)
			return
		}
		token = strings.TrimPrefix(token, Bearer)

		tokenHash := HashToken(token, []byte(config.App.HashSecret))
		retrievedToken, err := db.Queries.GetSessionToken(r.Context(), tokenHash)
		if err != nil {
			logger.Debug("Invalid token provided")
			http.Error(w, "Invalid token provided", http.StatusUnauthorized)
			return
		}

		uID := retrievedToken.UserID
		sID := retrievedToken.ID
		logger.Debug("Authenticated user", "userID", uID.String())
		ctx := context.WithValue(r.Context(), userID, uID)
		ctx = context.WithValue(ctx, sessionID, sID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserID(ctx context.Context) uuid.UUID {
	var id uuid.UUID

	id, ok := ctx.Value(userID).(uuid.UUID)
	if !ok {
		logger.Warn("Failed to fetch userID")
		return uuid.Nil
	}
	return id
}

func SessionID(ctx context.Context) uuid.UUID {
	var id uuid.UUID

	id, ok := ctx.Value(sessionID).(uuid.UUID)
	if !ok {
		logger.Warn("Failed to fetch sessionID")
		return uuid.Nil
	}
	return id
}

func AsUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, userID, id)
}
