package app

import (
	"context"
	"os"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func createTracer(ctx context.Context) (trace.Tracer, error) {
	uptrace.ConfigureOpentelemetry(
		uptrace.WithServiceName("vocabulary-leveling"),
		uptrace.WithServiceVersion("1.0.0"),
	)

	exporter, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint("otlp.uptrace.dev"),
		otlptracehttp.WithHeaders(map[string]string{
			"uptrace-dsn": os.Getenv("UPTRACE_DSN"),
		}),
		otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
	)
	if err != nil {
		return nil, err
	}

	bsp := sdktrace.NewBatchSpanProcessor(exporter,
		sdktrace.WithMaxQueueSize(10_000),
		sdktrace.WithMaxExportBatchSize(10_000))

	resource, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			attribute.String("service.name", "vocabulary-leveling"),
			attribute.String("service.version", "1.0.0"),
		))
	if err != nil {
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(resource),
		sdktrace.WithIDGenerator(xray.NewIDGenerator()),
	)
	tracerProvider.RegisterSpanProcessor(bsp)

	otel.SetTracerProvider(tracerProvider)
	tracer := otel.Tracer("vocubulary-leveling")

	return tracer, nil
}
