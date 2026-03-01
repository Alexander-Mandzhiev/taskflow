package metric

import (
	"testing"
)

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
