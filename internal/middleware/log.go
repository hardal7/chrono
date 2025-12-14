package middleware

import (
	"net/http"
	"time"

	logger "github.com/hardal7/chrono/internal/util"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Debug(r.Method + " " + r.URL.Path + " " + time.Since(start).String())
	})
}
