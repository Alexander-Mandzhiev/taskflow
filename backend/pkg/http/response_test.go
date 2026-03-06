package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	ctx := context.Background()
	w := httptest.NewRecorder()
	body := map[string]string{"key": "value"}
	WriteJSON(ctx, w, http.StatusCreated, body)
	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want 201", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	var decoded map[string]string
	if err := json.NewDecoder(w.Body).Decode(&decoded); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if decoded["key"] != "value" {
		t.Errorf("body = %v", decoded)
	}
}

func TestWriteJSON_NilBody(t *testing.T) {
	ctx := context.Background()
	w := httptest.NewRecorder()
	WriteJSON(ctx, w, http.StatusOK, nil)
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	if w.Body.Len() != 0 {
		t.Errorf("body should be empty, got %q", w.Body.String())
	}
}
