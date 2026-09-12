package telemetry

import (
	"context"
	"errors"
	"fmt"

	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func InitOTel(ctx context.Context) (func(context.Context), error) {
	endpoint := config.App.OTelEndpoint
	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName("api"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create resource: %w", err)
	}

	traceExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
	)
	if err != nil {
		return nil, fmt.Errorf("create trace exporter: %w", err)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithResource(res),
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
		return nil, fmt.Errorf("create metric exporter: %w", err)
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
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

	tracer := otel.Tracer("startup")
	_, span := tracer.Start(ctx, "startup-span")
	span.SetAttributes(
		attribute.String("startup", "healthy"),
	)
	span.End()

	logger.Info("Initialized OTel tracing")
	return func(ctx context.Context) {
		var errs []error

		if err := meterProvider.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}

		if err := tracerProvider.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}

		logger.Error(errors.Join(errs...).Error())
	}, nil
}
