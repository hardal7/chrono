package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/trace"
)

const collectorURL = ""

func InitTracer(ctx context.Context) (func(context.Context) error, error) {
	exporter, err := otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(collectorURL),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create exporter: %w", err)
	}

	// Configure OpenTelemetry tracer provider
	otel.SetTracerProvider(
		trace.NewTracerProvider(
			trace.WithSampler(trace.AlwaysSample()), // capture every trace without sampling, not advised in prod.
			trace.WithBatcher(exporter),             // batch spans for efficient delivery
		),
	)

	return exporter.Shutdown, err
}
