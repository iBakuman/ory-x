package otelx

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ITracer is the creator of Spans.
//
// Warning: Methods may be added to this interface in minor releases. See
// package documentation on API implementation for information on how to set
// default behavior for unimplemented methods.
type ITracer interface {

	// Start creates a span and a context.Context containing the newly-created span.
	//
	// If the context.Context provided in `ctx` contains a Span then the newly-created
	// Span will be a child of that span, otherwise it will be a root span. This behavior
	// can be overridden by providing `WithNewRoot()` as a SpanOption, causing the
	// newly-created Span to be a root span even if `ctx` contains a Span.
	//
	// When creating a Span it is recommended to provide all known span attributes using
	// the `WithAttributes()` SpanOption as samplers will only have access to the
	// attributes provided when a Span is created.
	//
	// Any Span that is created MUST also be ended. This is the responsibility of the user.
	// Implementations of this API may leak memory or other resources if Spans are not ended.
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, Span)
}

type TracerProvider interface {
	// Tracer returns a unique Tracer scoped to be used by instrumentation code
	// to trace computational workflows. The scope and identity of that
	// instrumentation code is uniquely defined by the name and options passed.
	//
	// The passed name needs to uniquely identify instrumentation code.
	// Therefore, it is recommended that name is the Go package name of the
	// library providing instrumentation (note: not the code being
	// instrumented). Instrumentation libraries can have multiple versions,
	// therefore, the WithInstrumentationVersion option should be used to
	// distinguish these different codebases. Additionally, instrumentation
	// libraries may sometimes use traces to communicate different domains of
	// workflow data (i.e. using spans to communicate workflow events only). If
	// this is the case, the WithScopeAttributes option should be used to
	// uniquely identify Tracers that handle the different domains of workflow
	// data.
	//
	// If the same name and options are passed multiple times, the same Tracer
	// will be returned (it is up to the implementation if this will be the
	// same underlying instance of that Tracer or not). It is not necessary to
	// call this multiple times with the same name and options to get an
	// up-to-date Tracer. All implementations will ensure any TracerProvider
	// configuration changes are propagated to all provided Tracers.
	//
	// If name is empty, then an implementation defined default name will be
	// used instead.
	//
	// This method is safe to call concurrently.
	Tracer(name string, options ...trace.TracerOption) ITracer
}

// Span is the individual component of a trace. It represents a single named
// and timed operation of a workflow that is traced. A Tracer is used to
// create a Span and it is then up to the operation the Span represents to
// properly end the Span when the operation itself ends.
//
// Warning: Methods may be added to this interface in minor releases. See
// package documentation on API implementation for information on how to set
// default behavior for unimplemented methods.
type Span interface {
	// End completes the Span. The Span is considered complete and ready to be
	// delivered through the rest of the telemetry pipeline after this method
	// is called. Therefore, updates to the Span are not allowed after this
	// method has been called.
	End(options ...trace.SpanEndOption)

	// AddEvent adds an event with the provided name and options.
	AddEvent(name string, options ...trace.EventOption)

	// IsRecording returns the recording state of the Span. It will return
	// true if the Span is active and events can be recorded.
	IsRecording() bool

	// RecordError will record err as an exception span event for this span. An
	// additional call to SetStatus is required if the Status of the Span should
	// be set to Error, as this method does not change the Span status. If this
	// span is not being recorded or err is nil then this method does nothing.
	RecordError(err error, options ...trace.EventOption)

	// SpanContext returns the SpanContext of the Span. The returned SpanContext
	// is usable even after the End method has been called for the Span.
	SpanContext() trace.SpanContext

	// SetStatus sets the status of the Span in the form of a code and a
	// description, provided the status hasn't already been set to a higher
	// value before (OK > Error > Unset). The description is only included in a
	// status when the code is for an error.
	SetStatus(code codes.Code, description string)

	// SetName sets the Span name.
	SetName(name string)

	// SetAttributes sets kv as attributes of the Span. If a key from kv
	// already exists for an attribute of the Span it will be overwritten with
	// the value contained in kv.
	SetAttributes(kv ...attribute.KeyValue)

	// TracerProvider returns a TracerProvider that can be used to generate
	// additional Spans on the same telemetry pipeline as the current Span.
	TracerProvider() trace.TracerProvider
}
