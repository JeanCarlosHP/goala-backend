package observability

import (
	"context"
	"strings"

	"github.com/jeancarloshp/calorieai/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace/noop"
)

func InitTracer(ctx context.Context, serviceName string, config *domain.Config) (*sdktrace.TracerProvider, error) {
	if !config.TracingEnabled || strings.TrimSpace(config.OtelCollectorURL) == "" {
		otel.SetTracerProvider(noop.NewTracerProvider())
		return nil, nil
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		otel.SetTracerProvider(noop.NewTracerProvider())
		return nil, err
	}

	otelCollector, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint(config.OtelCollectorURL), otlptracehttp.WithInsecure())
	if err != nil {
		otel.SetTracerProvider(noop.NewTracerProvider())
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(otelCollector),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}
