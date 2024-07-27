package otelx

import (
	"context"

	"github.com/theplant/appkit/logtracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type noopTracerProvider struct{}

// Tracer returns an OpenTelemetry Tracer that does not record any telemetry.
func (noopTracerProvider) Tracer(string, ...trace.TracerOption) ITracer {
	return noopTracer{}
}

type noopTracer struct{}

func (l noopTracer) Start(
	ctx context.Context, spanName string, opts ...trace.SpanStartOption,
) (context.Context, Span) {
	ctx, _ = logtracing.StartSpan(ctx, spanName)
	return ctx, logSpan{ctx: ctx}
}

// Span is an OpenTelemetry No-Op Span.
type noopSpan struct {
	sc trace.SpanContext
}

// SpanContext returns an empty span context.
func (s noopSpan) SpanContext() trace.SpanContext { return s.sc }

// IsRecording always returns false.
func (noopSpan) IsRecording() bool { return false }

// SetStatus does nothing.
func (noopSpan) SetStatus(codes.Code, string) {}

// SetAttributes does nothing.
func (noopSpan) SetAttributes(...attribute.KeyValue) {}

// End does nothing.
func (noopSpan) End(...trace.SpanEndOption) {}

// RecordError does nothing.
func (noopSpan) RecordError(error, ...trace.EventOption) {}

// AddEvent does nothing.
func (noopSpan) AddEvent(string, ...trace.EventOption) {}

// SetName does nothing.
func (noopSpan) SetName(string) {}

// TracerProvider returns a No-Op TracerProvider.
func (noopSpan) TracerProvider() trace.TracerProvider {
	panic("implement me")
}
