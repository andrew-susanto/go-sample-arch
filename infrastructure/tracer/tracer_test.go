package tracer

import (
	// golang package
	"context"
	"os"
	"testing"

	// external package
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/mock/gomock"
)

func TestTracer_NewExporter(t *testing.T) {
	tests := []struct {
		name string
		want func() sdktrace.SpanExporter
	}{
		{
			name: "test_return_new_exporter",
			want: func() sdktrace.SpanExporter {
				tracer, _ := stdouttrace.New(stdouttrace.WithWriter(os.Stdout), stdouttrace.WithPrettyPrint())
				return tracer
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			got := newExporter()
			assert.Equal(t, tt.want(), got)
		})
	}
}

func TestTracer_NewTraceProvider(t *testing.T) {
	type args struct {
		exp         sdktrace.SpanExporter
		serviceName string
		environment string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "when_given_exporter_and_service_name_then_expect_return_trace_provider",
			args: args{
				exp:         newExporter(),
				serviceName: "test",
				environment: "test",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			got := newTraceProvider(tt.args.exp, tt.args.serviceName, tt.args.environment)
			assert.NotNil(t, got)
		})
	}
}

func TestTracer_InitTracer(t *testing.T) {
	type args struct {
		serviceName string
		environment string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "when_given_service_name_expect_return_trace_provider",
			args: args{
				serviceName: "test",
				environment: "test",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			got := InitTracer(tt.args.serviceName, tt.args.environment)
			assert.NotNil(t, got)
		})
	}
}

func TestTracer_Start(t *testing.T) {
	type args struct {
		ctx      context.Context
		spanName string
	}
	tests := []struct {
		name string
		args args
		mock func()
	}{
		{
			name: "when_tracer_nil_expect_return_context_and_span",
			args: args{
				ctx:      context.Background(),
				spanName: "test_span",
			},
			mock: func() {
				tracer = nil
			},
		},
		{
			name: "when_tracer_not_nil_expect_return_context_and_span",
			args: args{
				ctx:      context.Background(),
				spanName: "test_span",
			},
			mock: func() {
				tracer = trace.NewNoopTracerProvider().Tracer("test")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tt.mock()

			gotContext, gotSpan := Start(tt.args.ctx, tt.args.spanName)
			assert.NotNil(t, gotContext)
			assert.NotNil(t, gotSpan)
		})
	}
}
