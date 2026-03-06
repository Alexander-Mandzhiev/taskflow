package tracing

import (
	"context"
	"testing"
	"time"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

func TestParseDuration(t *testing.T) {
	defaultVal := 5 * time.Second

	t.Run("valid", func(t *testing.T) {
		got := ParseDuration("10s", defaultVal)
		if got != 10*time.Second {
			t.Errorf("ParseDuration(\"10s\", 5s) = %v, want 10s", got)
		}
	})

	t.Run("invalid returns default", func(t *testing.T) {
		got := ParseDuration("invalid", defaultVal)
		if got != defaultVal {
			t.Errorf("ParseDuration(\"invalid\", 5s) = %v, want 5s", got)
		}
	})

	t.Run("empty returns default", func(t *testing.T) {
		got := ParseDuration("", defaultVal)
		if got != defaultVal {
			t.Errorf("ParseDuration(\"\", 5s) = %v, want 5s", got)
		}
	})
}

func TestNewWithLogger(t *testing.T) {
	t.Run("nil logger uses noop", func(t *testing.T) {
		tr := NewWithLogger(nil)
		if tr == nil {
			t.Fatal("NewWithLogger(nil) = nil")
		}
	})

	t.Run("with logger", func(t *testing.T) {
		tr := NewWithLogger(&logger.NoopLogger{})
		if tr == nil {
			t.Fatal("NewWithLogger(noop) = nil")
		}
	})
}

func TestTracing_SetLogger(t *testing.T) {
	tr := NewWithLogger(&logger.NoopLogger{})
	tr.SetLogger(nil) // не паникует и не меняет логгер
	tr.SetLogger(&logger.NoopLogger{})
}

func TestTracing_Init_Disabled(t *testing.T) {
	tr := NewWithLogger(&logger.NoopLogger{})
	ctx := context.Background()

	err := tr.Init(ctx, WithEnable(false))
	if err != nil {
		t.Errorf("Init(WithEnable(false)) = %v", err)
	}
}

func TestTracing_GetTracerProvider(t *testing.T) {
	tr := NewWithLogger(&logger.NoopLogger{})
	// до Init tracerProvider == nil — возвращается новый провайдер
	provider := tr.GetTracerProvider()
	if provider == nil {
		t.Error("GetTracerProvider() = nil")
	}
}

func TestTracing_Shutdown_NilProvider(t *testing.T) {
	tr := NewWithLogger(&logger.NoopLogger{})
	ctx := context.Background()

	err := tr.Shutdown(ctx, time.Second)
	if err != nil {
		t.Errorf("Shutdown when provider nil = %v", err)
	}
}

func TestTracing_Reinit(t *testing.T) {
	tr := NewWithLogger(&logger.NoopLogger{})
	ctx := context.Background()

	_ = tr.Init(ctx, WithEnable(false))
	err := tr.Reinit(ctx, WithEnable(false))
	if err != nil {
		t.Errorf("Reinit(WithEnable(false)) = %v", err)
	}
}
