package ctxkey

import (
	"context"
	"testing"
)

// Проверяем, что константы имеют ожидаемые значения и работают как ключи context.
func TestKey_constantsAsContextKeys(t *testing.T) {
	ctx := context.Background()

	ctx = context.WithValue(ctx, TraceID, "trace-123")
	ctx = context.WithValue(ctx, RequestID, "req-456")
	ctx = context.WithValue(ctx, UserID, "user-789")
	ctx = context.WithValue(ctx, SessionID, "sess-abc")

	if v, ok := ctx.Value(TraceID).(string); !ok || v != "trace-123" {
		t.Errorf("TraceID: got %q (ok=%v), want trace-123", v, ok)
	}
	if v, ok := ctx.Value(RequestID).(string); !ok || v != "req-456" {
		t.Errorf("RequestID: got %q (ok=%v), want req-456", v, ok)
	}
	if v, ok := ctx.Value(UserID).(string); !ok || v != "user-789" {
		t.Errorf("UserID: got %q (ok=%v), want user-789", v, ok)
	}
	if v, ok := ctx.Value(SessionID).(string); !ok || v != "sess-abc" {
		t.Errorf("SessionID: got %q (ok=%v), want sess-abc", v, ok)
	}
}

func TestKey_stringValues(t *testing.T) {
	want := map[Key]string{
		TraceID:   "trace_id",
		RequestID: "request_id",
		UserID:    "user_id",
		SessionID: "session_id",
	}
	for k, v := range want {
		if string(k) != v {
			t.Errorf("Key %q: string value = %q, want %q", k, string(k), v)
		}
	}
}
