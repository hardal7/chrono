package telemetry

import (
	"context"
	"errors"
	"fmt"

	"github.com/hardal7/chrono/internal/util/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

func InitOTel(ctx context.Context) (func(context.Context) error, error) {
	endpoint := config.App.OTelEndpoint

	traceExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create trace exporter: %w", err)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(traceExporter),
	)
	otel.SetTracerProvider(tracerProvider)

	metricExporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(endpoint),
	)
	if err != nil {
		_ = traceExporter.Shutdown(ctx)
		return nil, fmt.Errorf("Failed to create metric exporter: %w", err)
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(
			metric.NewPeriodicReader(metricExporter),
		),
	)
	otel.SetMeterProvider(meterProvider)
	if err := initHTTPMetrics(); err != nil {
		_ = meterProvider.Shutdown(ctx)
		_ = tracerProvider.Shutdown(ctx)

		return nil, fmt.Errorf("failed to initialize HTTP metrics: %w", err)
	}

	return func(ctx context.Context) error {
		var errs []error

		if err := meterProvider.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}

		if err := tracerProvider.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}

		return errors.Join(errs...)
	}, nil
}
