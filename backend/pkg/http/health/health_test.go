package healthhttp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	Handler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type = %q", w.Header().Get("Content-Type"))
	}
	var resp Response
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Status != "ok" {
		t.Errorf("status = %q, want ok", resp.Status)
	}
}

func TestLiveHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/live", nil)
	w := httptest.NewRecorder()
	LiveHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("LiveHandler status = %d", w.Code)
	}
}

func TestReadyHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()
	ReadyHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("ReadyHandler status = %d", w.Code)
	}
}

func TestStartHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/start", nil)
	w := httptest.NewRecorder()
	StartHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("StartHandler status = %d", w.Code)
	}
}
