package telemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	httpRequestCount    metric.Int64Counter
	httpRequestDuration metric.Float64Histogram
	httpActiveRequests  metric.Int64UpDownCounter
)

func RecordHTTPActiveRequest(
	ctx context.Context,
	method string,
	route string,
) {
	attrs := metric.WithAttributes(
		attribute.String("http.request.method", method),
		attribute.String("http.route", route),
	)

	httpActiveRequests.Add(ctx, 1, attrs)
}

func RecordHTTPResponse(
	ctx context.Context,
	method string,
	route string,
	status int,
	duration time.Duration,
) {
	attrs := metric.WithAttributes(
		attribute.String("http.request.method", method),
		attribute.String("http.route", route),
		attribute.Int("http.response.status_code", status),
	)

	httpRequestCount.Add(ctx, 1, attrs)

	httpRequestDuration.Record(
		ctx,
		duration.Seconds(),
		attrs,
	)

	httpActiveRequests.Add(
		ctx,
		-1,
		metric.WithAttributes(
			attribute.String("http.request.method", method),
			attribute.String("http.route", route),
		),
	)
}

func initHTTPMetrics() error {
	meter := otel.Meter("api/http")

	var err error

	httpRequestCount, err = meter.Int64Counter(
		"http.server.request.count",
		metric.WithDescription("Number of HTTP requests received"),
	)
	if err != nil {
		return err
	}

	httpRequestDuration, err = meter.Float64Histogram(
		"http.server.request.duration",
		metric.WithUnit("s"),
		metric.WithDescription("Duration of HTTP server requests"),
	)
	if err != nil {
		return err
	}

	httpActiveRequests, err = meter.Int64UpDownCounter(
		"http.server.active_requests",
		metric.WithDescription("Number of currently active HTTP requests"),
	)
	if err != nil {
		return err
	}

	return nil
}
