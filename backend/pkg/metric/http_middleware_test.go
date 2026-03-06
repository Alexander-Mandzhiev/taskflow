package metric

import (
	"errors"
	"net/http"
	"testing"
)

func TestShouldRecordMetrics(t *testing.T) {
	t.Parallel()
	tests := []struct {
		path string
		want bool
	}{
		{"", false},
		{"/", false},
		{"/api", true},
		{"/api/v1/users", true},
		{"/health", true},
		{"/healthz", true},
		{"/live", true},
		{"/ready", true},
		{"/start", true},
		{"/other", false},
		{"/metrics", false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.path, func(t *testing.T) {
			t.Parallel()
			if got := shouldRecordMetrics(tt.path); got != tt.want {
				t.Errorf("shouldRecordMetrics(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestStatusWriter_WriteHeader(t *testing.T) {
	rw := &mockResponseWriter{}
	sw := &statusWriter{ResponseWriter: rw, status: 0, written: false}
	sw.WriteHeader(http.StatusCreated)
	if sw.status != http.StatusCreated {
		t.Errorf("status = %d, want 201", sw.status)
	}
	if !sw.written {
		t.Error("written should be true")
	}
	sw.WriteHeader(http.StatusOK) // повторный вызов игнорируется
	if sw.status != http.StatusCreated {
		t.Errorf("second WriteHeader changed status to %d", sw.status)
	}
}

func TestStatusWriter_Write_SetsStatusOK(t *testing.T) {
	rw := &mockResponseWriter{}
	sw := &statusWriter{ResponseWriter: rw, status: 0, written: false}
	_, _ = sw.Write([]byte("ok"))
	if sw.status != http.StatusOK {
		t.Errorf("Write without WriteHeader: status = %d, want 200", sw.status)
	}
	if !sw.written {
		t.Error("written should be true after Write")
	}
}

func TestStatusWriter_Hijack_ErrNotHijacker(t *testing.T) {
	rw := &mockResponseWriter{} // не реализует Hijacker
	sw := &statusWriter{ResponseWriter: rw}
	_, _, err := sw.Hijack()
	if !errors.Is(err, ErrNotHijacker) {
		t.Errorf("Hijack() = %v, want ErrNotHijacker", err)
	}
}

func TestStatusWriter_Push_ErrNotPusher(t *testing.T) {
	rw := &mockResponseWriter{} // не реализует Pusher
	sw := &statusWriter{ResponseWriter: rw}
	err := sw.Push("/", nil)
	if !errors.Is(err, ErrNotPusher) {
		t.Errorf("Push() = %v, want ErrNotPusher", err)
	}
}

type mockResponseWriter struct {
	header http.Header
	code   int
	body   []byte
}

func (m *mockResponseWriter) Header() http.Header {
	if m.header == nil {
		m.header = make(http.Header)
	}
	return m.header
}

func (m *mockResponseWriter) Write(b []byte) (int, error) {
	m.body = append(m.body, b...)
	return len(b), nil
}
func (m *mockResponseWriter) WriteHeader(code int) { m.code = code }

func TestNormalizePath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", "/"},
		{"root", "/", "/"},
		{"trailing_slash", "/api/v1/", "/api/v1"},
		{"double_slashes", "//api//v1///health", "/api/v1/health"},
		{"uuid_segment", "/api/v1/teams/e5fd738d-dfef-4458-8379-13039bbd6a63/invite", "/api/v1/teams/{id}/invite"},
		{"int_segment", "/api/v1/users/123/profile", "/api/v1/users/{id}/profile"},
		{"task_uuid", "/api/v1/tasks/e5fd738d-dfef-4458-8379-13039bbd6a63/history", "/api/v1/tasks/{id}/history"},
		{"hex24_segment", "/api/v1/items/507f1f77bcf86cd799439011", "/api/v1/items/{id}"},
		{"hex24_before_int", "/api/v1/items/123456789012345678901234", "/api/v1/items/{id}"}, // hex24 должен обрабатываться ДО int
		{"many_slashes", "////api////v1////users////1////", "/api/v1/users/{id}"},
		{"no_leading_slash", "api/v1/users/1", "/api/v1/users/{id}"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := normalizePath(tt.in); got != tt.want {
				t.Fatalf("normalizePath(%q)=%q; want %q", tt.in, got, tt.want)
			}
		})
	}
}
