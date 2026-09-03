package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/hardal7/chrono/internal/util/requestctx"
)

const (
	AuthCookie = "Authorization"
	Bearer     = "Bearer "
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Authenticating user")

		cookie, err := r.Cookie(AuthCookie)
		if cookie == nil {
			logger.Debug("No token provided")
			http.Error(w, "No token provided", http.StatusUnauthorized)
			return
		} else if err != nil {
			logger.Debug("Invalid cookie")
			http.Error(w, "Invalid cookie", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(cookie.Value, Bearer)

		tokenHash := HashToken(token, []byte(config.App.HashSecret))
		retrievedToken, err := db.Queries.GetSessionToken(r.Context(), tokenHash)
		if err != nil {
			logger.Debug("Invalid token provided")
			http.Error(w, "Invalid token provided", http.StatusUnauthorized)
			return
		}

		userID := retrievedToken.UserID
		sessionID := retrievedToken.ID
		logger.Debug("Authenticated user", "userID", userID.String())
		ctx := context.WithValue(r.Context(), requestctx.UserID, userID)
		ctx = context.WithValue(ctx, requestctx.SessionID, sessionID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
