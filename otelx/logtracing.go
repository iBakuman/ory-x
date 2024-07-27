package otelx

import (
	"context"

	"github.com/theplant/appkit/logtracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type LogTracingTracer struct {
}

func (l LogTracingTracer) Tracer() {
}

func (l LogTracingTracer) Start(
	ctx context.Context, spanName string, opts ...trace.SpanStartOption,
) (context.Context, Span) {
	ctx, _ = logtracing.StartSpan(ctx, spanName)
	return ctx, logSpan{ctx: ctx}
}

func SetupLogtracing() *LogTracingTracer {
	return &LogTracingTracer{}
}

type logSpan struct {
	ctx context.Context
}

func (l logSpan) End(options ...trace.SpanEndOption) {
	if s := logtracing.SpanFromContext(l.ctx); s != nil {
		s.End()
	}
}

func (l logSpan) AddEvent(name string, options ...trace.EventOption) {
	if s := logtracing.SpanFromContext(l.ctx); s != nil {
		logtracing.AppendSpanKVs(l.ctx, "kratos.event", name)
	}
}

func (l logSpan) IsRecording() bool {
	if s := logtracing.SpanFromContext(l.ctx); s != nil {
		return s.IsRecording()
	} else {
		return false
	}
}

func (l logSpan) RecordError(err error, options ...trace.EventOption) {
	if s := logtracing.SpanFromContext(l.ctx); s != nil {
		s.RecordError(err)
	}
}

func (l logSpan) SpanContext() trace.SpanContext {
	return trace.SpanContext{}
}

func (l logSpan) SetStatus(code codes.Code, description string) {
}

func (l logSpan) SetName(name string) {
}

func (l logSpan) SetAttributes(kv ...attribute.KeyValue) {
	if s := logtracing.SpanFromContext(l.ctx); s != nil {
		for _, k := range kv {
			logtracing.AppendSpanKVs(l.ctx, k.Key, k.Value)
		}
	}
}

func (l logSpan) TracerProvider() trace.TracerProvider {
	panic("implement me")
}
