package middleware

import (
	"net/http"
	"strconv"
	"time"

	logger "github.com/hardal7/chrono/internal/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(ww, r)

		duration := time.Since(start)
		method := r.Method
		endpoint := r.URL.Path
		status := ww.status

		logger.Debug(strconv.Itoa(status) + " " + method + " " + endpoint + " " + duration.String())
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
