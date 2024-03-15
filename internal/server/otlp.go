package server

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func initTracer() {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("auth-api"),
			semconv.ServiceVersion("v0.0.1"),
		),
	)
	if err != nil {
		panic(err)
	}
	// Create the exporter - let's use a stdout exporter

	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(os.Getenv("OTLP_ENDPOINT")),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithCompression(otlptracehttp.NoCompression))
	traceExporter, err := otlptrace.New(ctx, client)
	if err != nil {
		panic(err)
	}

	// Configure the trace provider
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter, trace.WithBatchTimeout(5*time.Millisecond)),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(traceProvider)

}
