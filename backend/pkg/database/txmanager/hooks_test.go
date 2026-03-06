package txmanager

import (
	"context"
	"testing"
)

func TestNewHookRegistry(t *testing.T) {
	r := NewHookRegistry()
	if r == nil {
		t.Fatal("NewHookRegistry() = nil")
	}
	hooks := r.GetHooks()
	if len(hooks) != 0 {
		t.Errorf("GetHooks() on new registry: len = %d, want 0", len(hooks))
	}
}

func TestHookRegistry_Register_GetHooks(t *testing.T) {
	r := NewHookRegistry()
	callOrder := []string{}
	hook1 := func(ctx context.Context) error {
		callOrder = append(callOrder, "a")
		return nil
	}
	hook2 := func(ctx context.Context) error {
		callOrder = append(callOrder, "b")
		return nil
	}

	r.Register("first", hook1)
	r.Register("second", hook2)
	hooks := r.GetHooks()
	if len(hooks) != 2 {
		t.Fatalf("GetHooks() len = %d, want 2", len(hooks))
	}
	ctx := context.Background()
	_ = hooks[0](ctx)
	_ = hooks[1](ctx)
	if callOrder[0] != "a" || callOrder[1] != "b" {
		t.Errorf("hook order = %v, want [a b]", callOrder)
	}
}

func TestHookRegistry_Register_skipsEmptyKeyAndNilHook(t *testing.T) {
	r := NewHookRegistry()
	r.Register("", func(context.Context) error { return nil })
	r.Register("key", nil)
	hooks := r.GetHooks()
	if len(hooks) != 0 {
		t.Errorf("Register empty key or nil hook should be skipped, got %d hooks", len(hooks))
	}
}

func TestHookRegistry_Register_overwritePreservesOrder(t *testing.T) {
	r := NewHookRegistry()
	r.Register("a", func(context.Context) error { return nil })
	r.Register("b", func(context.Context) error { return nil })
	r.Register("a", func(context.Context) error { return nil }) // перезапись
	hooks := r.GetHooks()
	if len(hooks) != 2 {
		t.Errorf("GetHooks() len = %d, want 2 (order preserved)", len(hooks))
	}
}

func TestWithHookRegistry_GetHookRegistry(t *testing.T) {
	ctx := context.Background()
	if GetHookRegistry(ctx) != nil {
		t.Error("GetHookRegistry(empty ctx) should be nil")
	}
	r := NewHookRegistry()
	ctx = WithHookRegistry(ctx, r)
	got := GetHookRegistry(ctx)
	if got != r {
		t.Errorf("GetHookRegistry(WithHookRegistry(ctx, r)) = %p, want %p", got, r)
	}
}
