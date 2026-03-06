package closer

import (
	"context"
	"errors"
	"testing"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

// Тесты без сигналов, чтобы не запускать handleSignals.

func TestNew(t *testing.T) {
	c := New()
	if c == nil {
		t.Fatal("New() = nil")
	}
}

func TestNewWithLogger_NilLogger(t *testing.T) {
	c := NewWithLogger(nil)
	if c == nil {
		t.Fatal("NewWithLogger(nil) = nil")
	}
}

func TestCloser_Add_CloseAll_Order(t *testing.T) {
	c := NewWithLogger(&logger.NoopLogger{})
	order := []int{}
	c.Add(
		func(context.Context) error { order = append(order, 1); return nil },
		func(context.Context) error { order = append(order, 2); return nil },
		func(context.Context) error { order = append(order, 3); return nil },
	)
	err := c.CloseAll(context.Background())
	if err != nil {
		t.Fatalf("CloseAll: %v", err)
	}
	// LIFO: 3, 2, 1
	if len(order) != 3 || order[0] != 3 || order[1] != 2 || order[2] != 1 {
		t.Errorf("CloseAll order = %v, want [3 2 1]", order)
	}
}

func TestCloser_CloseAll_ReturnsFirstError(t *testing.T) {
	c := NewWithLogger(&logger.NoopLogger{})
	wantErr := errors.New("close failed")
	c.Add(func(context.Context) error { return wantErr })
	err := c.CloseAll(context.Background())
	if !errors.Is(err, wantErr) {
		t.Errorf("CloseAll = %v, want %v", err, wantErr)
	}
}

func TestCloser_CloseAll_Idempotent(t *testing.T) {
	c := NewWithLogger(&logger.NoopLogger{})
	calls := 0
	c.Add(func(context.Context) error { calls++; return nil })
	_ = c.CloseAll(context.Background())
	_ = c.CloseAll(context.Background())
	if calls != 1 {
		t.Errorf("close fn called %d times, want 1", calls)
	}
}

func TestCloser_CloseAll_Empty(t *testing.T) {
	c := NewWithLogger(&logger.NoopLogger{})
	err := c.CloseAll(context.Background())
	if err != nil {
		t.Errorf("CloseAll(empty) = %v", err)
	}
}

func TestCloser_CloseAll_ContextCancelled(t *testing.T) {
	c := NewWithLogger(&logger.NoopLogger{})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c.Add(func(context.Context) error { return nil })
	err := c.CloseAll(ctx)
	if err == nil {
		t.Fatal("CloseAll(cancelled ctx) expected error")
	}
}

func TestCloser_SetLogger(t *testing.T) {
	c := NewWithLogger(&logger.NoopLogger{})
	c.SetLogger(nil) // не паникует
	c.SetLogger(&logger.NoopLogger{})
}
