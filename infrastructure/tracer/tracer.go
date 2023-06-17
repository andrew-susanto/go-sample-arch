package tracer

import (
	// golang package
	"context"
	"os"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"

	// external package
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

// newExporter creates a stdout exporter to be used with a trace provider.
func newExporter() sdktrace.SpanExporter {
	exporter, err := stdouttrace.New(stdouttrace.WithWriter(os.Stdout), stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Error(err, nil, "stdouttrace.New() got error - newExporter")
		return nil
	}

	return exporter
}

// newTraceProvider creates a new trace provider instance.
func newTraceProvider(exp sdktrace.SpanExporter, serviceName string, environment string) *sdktrace.TracerProvider {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.DeploymentEnvironment(environment),
		),
	)
	if err != nil {
		log.Fatal(err, nil, "resource.Megre() got error - newTraceProvider")
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
}

// InitTracer creates a new trace provider instance and registers it as global trace provider.
func InitTracer(serviceName string, environment string) *sdktrace.TracerProvider {
	exp := newExporter()
	tp := newTraceProvider(exp, serviceName, environment)

	otel.SetTracerProvider(tp)

	// Finally, set the tracer that can be used for this package.
	tracer = tp.Tracer(serviceName)

	return tp
}

// Start starts a span with given span name and returns the span.
func Start(ctx context.Context, spanName string) (context.Context, trace.Span) {
	if tracer != nil {
		return tracer.Start(ctx, spanName)
	}

	return sdktrace.NewTracerProvider().Tracer("default").Start(ctx, spanName)
}
