package tracing

import (
	"context"
	"testing"
)

func TestTraceIDFromContext_Empty(t *testing.T) {
	ctx := context.Background()
	got := TraceIDFromContext(ctx)
	if got != "" {
		t.Errorf("TraceIDFromContext(empty ctx) = %q, want \"\"", got)
	}
}

func TestStartSpan_TraceIDInContext(t *testing.T) {
	ctx := context.Background()
	newCtx, span := StartSpan(ctx, "test")
	if span == nil {
		t.Fatal("StartSpan returned nil span")
	}
	defer span.End()

	traceID := TraceIDFromContext(newCtx)
	if traceID == "" {
		t.Error("TraceIDFromContext after StartSpan should be non-empty")
	}
}

func TestSpanFromContext_Empty(t *testing.T) {
	ctx := context.Background()
	span := SpanFromContext(ctx)
	if span == nil {
		t.Fatal("SpanFromContext(empty) = nil")
	}
	// у пустого контекста спан noop, IsValid() может быть false
	_ = span
}

func TestSpanFromContext_AfterStartSpan(t *testing.T) {
	ctx := context.Background()
	newCtx, _ := StartSpan(ctx, "child")
	span := SpanFromContext(newCtx)
	if span == nil {
		t.Fatal("SpanFromContext(after StartSpan) = nil")
	}
	if !span.SpanContext().IsValid() {
		t.Error("span from context after StartSpan should be valid")
	}
}
