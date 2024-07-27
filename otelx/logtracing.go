package otelx

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type LogTracingTracer struct {
}

func (l LogTracingTracer) tracer() {
	// TODO implement me
	panic("implement me")
}

func (l LogTracingTracer) Start(
	ctx context.Context, spanName string, opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	// TODO implement me
	panic("implement me")
}

func SetupLogtracing() *LogTracingTracer {
	return &LogTracingTracer{}
}
