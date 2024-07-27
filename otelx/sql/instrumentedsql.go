// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"context"
	"database/sql/driver"

	"github.com/luna-duclos/instrumentedsql"
	"github.com/pkg/errors"
	"github.com/theplant/appkit/logtracing"
)

const tracingComponent = "github.com/ory/x/otelx/sql"

type (
	tracer struct{}
	span   struct {
		ctx context.Context
	}
)

var (
	_ instrumentedsql.Tracer = tracer{}
	_ instrumentedsql.Span   = span{}
)

func NewTracer() instrumentedsql.Tracer { return tracer{} }

// GetSpan returns a span
func (tracer) GetSpan(ctx context.Context) instrumentedsql.Span {
	ctx, _ = logtracing.StartSpan(ctx, tracingComponent)
	return span{ctx: ctx}
}

func (s span) NewChild(name string) instrumentedsql.Span {
	nCtx, _ := logtracing.StartSpan(s.ctx, name)
	return span{ctx: nCtx}
}

func (s span) SetLabel(k, v string) {
	logtracing.AppendSpanKVs(s.ctx, k, v)
}

func (s span) SetError(err error) {
	if err == nil || errors.Is(err, driver.ErrSkip) {
		return
	}
	logtracing.AppendSpanKVs(s.ctx, "error", err.Error())
}

func (s span) Finish() {
	t := logtracing.SpanFromContext(s.ctx)
	if t != nil {
		t.End()
	}
}
