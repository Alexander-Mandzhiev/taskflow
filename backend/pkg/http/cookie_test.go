package http

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSetCookie(t *testing.T) {
	w := httptest.NewRecorder()
	SetCookie(w, "session_id", "abc123", 3600, false, "", true)
	header := w.Header().Get("Set-Cookie")
	if header == "" {
		t.Fatal("Set-Cookie header should be set")
	}
	if !strings.Contains(header, "session_id=abc123") {
		t.Errorf("Set-Cookie should contain name=value: %s", header)
	}
	if !strings.Contains(header, "HttpOnly") {
		t.Error("cookie should be HttpOnly")
	}
}

func TestDeleteCookie(t *testing.T) {
	w := httptest.NewRecorder()
	DeleteCookie(w, "session_id", false, "")
	header := w.Header().Get("Set-Cookie")
	if header == "" {
		t.Fatal("Set-Cookie header should be set for delete")
	}
	// Удаление cookie задаётся Max-Age=-1 или Expires в прошлом; сериализация может дать Max-Age=0
	if !strings.Contains(header, "Max-Age=") && !strings.Contains(header, "Expires=") {
		t.Errorf("delete cookie should set expiry: %s", header)
	}
}
