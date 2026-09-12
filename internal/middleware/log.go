package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"uuid"

	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/hardal7/chrono/internal/util/requestctx"
	"github.com/hardal7/chrono/internal/util/telemetry"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Err(err).Debug("Failed to read body")
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(body))

		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), requestctx.RequestID, requestID)

		logger.Trace("Received Request")

		contentType := r.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "application/json") ||
			strings.HasPrefix(contentType, "text/plain") {
			// TODO: Do not print sensitive information
			logger.With("requestID", requestID).Debug(string(body))
		} else {
			logger.With("Content-Type", contentType, "bytes", len(body)).
				Debug("Request body omitted")
		}

		address := r.Header.Get("X-Forwarded-For")
		ctx = context.WithValue(ctx, requestctx.IP, address)

		route := r.URL.Path

		telemetry.RecordHTTPActiveRequest(ctx, r.Method, route)

		start := time.Now()

		ww := &statusWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(ww, r.WithContext(ctx))

		duration := time.Since(start)
		status := ww.status

		telemetry.RecordHTTPResponse(
			ctx,
			r.Method,
			route,
			status,
			duration,
		)

		logger.Info(fmt.Sprintf(
			"%d %s %s %s %s %s",
			status,
			r.Method,
			r.URL.Path,
			address,
			duration,
			requestID,
		))
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
	w.status = status
}
