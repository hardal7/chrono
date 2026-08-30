package middleware

import (
	"net/http"

	"github.com/hardal7/chrono/internal/auth"
	"github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/util/logger"
)

func Activity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Tracking user activity")

		userID := auth.UserID(r.Context())
		err := db.Queries.UpdateUserActivity(r.Context(), userID)
		if err != nil {
			logger.Warn("Failed to update user activity")
		}

		next.ServeHTTP(w, r)
	})
}
