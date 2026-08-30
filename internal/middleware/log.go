package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/hardal7/chrono/internal/util/requestctx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Debug("Failed to read body", "error", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(body))

		logger.Trace("Received Request")
		contentType := r.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "application/json") || strings.HasPrefix(contentType, "text/plain") {
			// TODO: Do not print sensitive information
			logger.Debug(string(body))
		} else {
			logger.Debug(fmt.Sprintf(
				"Request body omitted (Content-Type: %s, %d bytes)",
				contentType,
				len(body),
			))
		}

		address := r.Header.Get("X-Forwarded-For")
		ctx := context.WithValue(r.Context(), requestctx.IP, address)

		requestID := uuid.New().String()
		ctx = context.WithValue(ctx, requestctx.RequestID, requestID)

		start := time.Now()
		ww := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(ww, r.WithContext(ctx))

		duration := time.Since(start)
		method := r.Method
		endpoint := r.URL.Path
		status := ww.status
		logger.Info(strconv.Itoa(status) + " " + method + " " + endpoint + " " + address + " " + duration.String() + " " + requestID)
		httpRequestsTotal.WithLabelValues(method, endpoint, http.StatusText(status)).Inc()
		httpRequestDuration.WithLabelValues(method, endpoint).Observe(float64(duration.Milliseconds()))
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

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "Total number of API requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_request_duration_milliseconds",
			Help:    "API request duration in milliseconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)
