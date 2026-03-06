package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHasAnyPrefix(t *testing.T) {
	tests := []struct {
		s        string
		prefixes []string
		want     bool
	}{
		{"/wp-admin", []string{"/wp-", "/wordpress"}, true},
		{"/wordpress", []string{"/wp-", "/wordpress"}, true},
		{"/api/v1", []string{"/wp-"}, false},
		{"", []string{"/wp-"}, false},
		{"/actuator/health", []string{"/actuator"}, true},
	}
	for _, tt := range tests {
		got := hasAnyPrefix(tt.s, tt.prefixes...)
		if got != tt.want {
			t.Errorf("hasAnyPrefix(%q, %v) = %v, want %v", tt.s, tt.prefixes, got, tt.want)
		}
	}
}

func TestRequestFirewallMiddleware_AllowsHealth(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := RequestFirewallMiddleware(next)
	paths := []string{"/health", "/healthz", "/live", "/ready", "/start"}
	for _, p := range paths {
		req := httptest.NewRequest(http.MethodGet, p, nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("path %q: status = %d, want 200", p, w.Code)
		}
	}
}

func TestRequestFirewallMiddleware_AllowsAPI(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := RequestFirewallMiddleware(next)
	paths := []string{"/api", "/api/v1", "/api/v1/login"}
	for _, p := range paths {
		req := httptest.NewRequest(http.MethodGet, p, nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("path %q: status = %d, want 200", p, w.Code)
		}
	}
}

func TestRequestFirewallMiddleware_BlocksScannerPaths(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := RequestFirewallMiddleware(next)
	paths := []string{"/.env", "/.git/config", "/wp-admin", "/actuator/health", "/server-status", "/api/../etc/passwd"}
	for _, p := range paths {
		req := httptest.NewRequest(http.MethodGet, p, nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Errorf("path %q: status = %d, want 404", p, w.Code)
		}
	}
}

func TestRequestFirewallMiddleware_BlocksUnknown(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := RequestFirewallMiddleware(next)
	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}
