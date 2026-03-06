package http_router

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http/middleware"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

func init() {
	logger.SetNopLogger()
}

func TestNewRouter_ReturnsNonNil(t *testing.T) {
	r := NewRouter(5*time.Second, nil)
	if r == nil {
		t.Fatal("router is nil")
	}
}

func TestNewRouter_HealthRequestSucceeds(t *testing.T) {
	r := NewRouter(5*time.Second, nil)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		data := []byte("ok")
		n, err := w.Write(data)
		if err != nil {
			t.Errorf("w.Write: %v", err)
		}
		if n != len(data) {
			t.Errorf("w.Write wrote %d bytes, want %d", n, len(data))
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /health status = %d, want 200", w.Code)
	}
	if body := w.Body.String(); body != "ok" {
		t.Errorf("body = %q, want \"ok\"", body)
	}
}

func TestNewRouter_APIPathSucceeds(t *testing.T) {
	r := NewRouter(5*time.Second, nil)

	r.Get("/api/v1/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /api/v1/ping status = %d, want 200", w.Code)
	}
}

func TestNewRouter_UnknownPathReturns404(t *testing.T) {
	r := NewRouter(5*time.Second, nil)

	r.Get("/api/v1/ok", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("GET /unknown status = %d, want 404", w.Code)
	}
}

func TestNewRouter_HealthPathsServed(t *testing.T) {
	r := NewRouter(5*time.Second, nil)

	for _, path := range []string{"/health", "/live", "/ready", "/start"} {
		r.Get(path, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	}

	for _, path := range []string{"/health", "/live", "/ready", "/start"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("GET %s status = %d, want 200", path, w.Code)
		}
	}
}

func TestNewRouter_WithCORS_AddsHeaders(t *testing.T) {
	origins := []string{"https://app.example.com"}
	methods := []string{"GET", "POST"}
	headers := []string{"Content-Type"}
	corsMw := middleware.CORSMiddleware(origins, methods, headers, nil, true, 600)
	r := NewRouter(5*time.Second, []func(http.Handler) http.Handler{corsMw})

	r.Get("/api/v1/me", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	req.Header.Set("Origin", "https://app.example.com")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "https://app.example.com" {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q",
			w.Header().Get("Access-Control-Allow-Origin"), "https://app.example.com")
	}
	if w.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Error("expected Access-Control-Allow-Credentials: true")
	}
}

func TestNewRouter_WithoutCORS_NoOriginHeader(t *testing.T) {
	r := NewRouter(5*time.Second, nil)

	r.Get("/api/v1/me", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	req.Header.Set("Origin", "https://evil.com")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Errorf("without CORS config, Access-Control-Allow-Origin should be empty, got %q",
			w.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestNewRouter_OptionsPreflight(t *testing.T) {
	origins := []string{"https://app.example.com"}
	methods := []string{"GET", "POST"}
	corsMw := middleware.CORSMiddleware(origins, methods, nil, nil, false, 0)
	r := NewRouter(5*time.Second, []func(http.Handler) http.Handler{corsMw})

	r.Get("/api/v1/me", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/me", nil)
	req.Header.Set("Origin", "https://app.example.com")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("OPTIONS status = %d, want 204", w.Code)
	}
}
