package tracing

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

// Init returns an instance of Jaeger Tracer.
func NewTracingService(ctx context.Context, service string) trace.Tracer {
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
	)
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		log.Fatal("creating OTLP trace exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(service),
		)),
	)

	return tp.Tracer(service)
}
