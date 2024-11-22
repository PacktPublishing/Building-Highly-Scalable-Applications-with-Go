package jaeger

import (
	"context"
	"fmt"

	"github.com/hannesdejager/utxo-tracker/internal/app/config"
	"github.com/hannesdejager/utxo-tracker/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func InitTracing(c config.Tracing, info domain.ServiceInstance) (
	*trace.TracerProvider,
	error,
) {
	exp, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpointURL(c.ExportEndpointURL),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to initialize Jaeger exporter: %w",
			err,
		)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(info.Name),
			semconv.ServiceNamespaceKey.String("utxo-tracker"),
			semconv.ServiceVersionKey.String(
				info.Version.CommitShortHash),
		)),
	)
	otel.SetTracerProvider(tp)

	return tp, nil
}
