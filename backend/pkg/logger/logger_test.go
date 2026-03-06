package logger

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

func TestWithIDs(t *testing.T) {
	ctx := context.Background()

	t.Run("both ids", func(t *testing.T) {
		ctx := WithIDs(ctx, "trace-1", "req-1")
		if got := TraceIDFrom(ctx); got != "trace-1" {
			t.Errorf("TraceIDFrom = %q, want trace-1", got)
		}
		if got := RequestIDFrom(ctx); got != "req-1" {
			t.Errorf("RequestIDFrom = %q, want req-1", got)
		}
	})

	t.Run("only requestID", func(t *testing.T) {
		ctx := context.Background()
		ctx = WithIDs(ctx, "", "req-only")
		if got := TraceIDFrom(ctx); got != "" {
			t.Errorf("TraceIDFrom = %q, want empty", got)
		}
		if got := RequestIDFrom(ctx); got != "req-only" {
			t.Errorf("RequestIDFrom = %q, want req-only", got)
		}
	})

	t.Run("only traceID", func(t *testing.T) {
		ctx := context.Background()
		ctx = WithIDs(ctx, "trace-only", "")
		if got := TraceIDFrom(ctx); got != "trace-only" {
			t.Errorf("TraceIDFrom = %q, want trace-only", got)
		}
		if got := RequestIDFrom(ctx); got != "" {
			t.Errorf("RequestIDFrom = %q, want empty", got)
		}
	})
}

func TestTraceIDFrom(t *testing.T) {
	if got := TraceIDFrom(context.Background()); got != "" {
		t.Errorf("TraceIDFrom(empty ctx) = %q, want \"\"", got)
	}
}

func TestRequestIDFrom(t *testing.T) {
	if got := RequestIDFrom(context.Background()); got != "" {
		t.Errorf("RequestIDFrom(empty ctx) = %q, want \"\"", got)
	}
}

func TestSetNopLogger_Logger(t *testing.T) {
	SetNopLogger()
	defer SetNopLogger() // оставляем nop для других тестов

	if Logger() == nil {
		t.Error("Logger() after SetNopLogger is nil")
	}
}

func TestSync_NoPanic(t *testing.T) {
	SetNopLogger()
	err := Sync()
	if err != nil {
		t.Errorf("Sync() with nop logger = %v", err)
	}
}

func TestWith_NoPanic(t *testing.T) {
	SetNopLogger()
	l := With(zap.String("key", "val"))
	if l == nil {
		t.Error("With(...) returned nil")
	}
}

func TestWithContext_NoPanic(t *testing.T) {
	SetNopLogger()
	ctx := WithIDs(context.Background(), "t", "r")
	l := WithContext(ctx)
	if l == nil {
		t.Error("WithContext(...) returned nil")
	}
}

// SetLevel не тестируем без полной Init: он использует глобальный level (zap.AtomicLevel),
// который инициализируется только в initLogger.

func TestDebugInfoWarnError_NoPanic(t *testing.T) {
	SetNopLogger()
	ctx := context.Background()
	Debug(ctx, "msg")
	Info(ctx, "msg")
	Warn(ctx, "msg")
	Error(ctx, "msg")
}
