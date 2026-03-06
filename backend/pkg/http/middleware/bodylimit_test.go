package middleware

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsBodyError_EOF(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := context.Background()
	handled := IsBodyError(ctx, w, io.EOF)
	if !handled {
		t.Error("IsBodyError(io.EOF) = false, want true")
	}
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestIsBodyError_MaxBytesError(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := context.Background()
	err := &http.MaxBytesError{Limit: 100}
	handled := IsBodyError(ctx, w, err)
	if !handled {
		t.Error("IsBodyError(MaxBytesError) = false, want true")
	}
	if w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("status = %d, want 413", w.Code)
	}
}

func TestIsBodyError_OtherError(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := context.Background()
	handled := IsBodyError(ctx, w, io.ErrClosedPipe)
	if handled {
		t.Error("IsBodyError(other) = true, want false")
	}
	// Для не-тела ошибки не пишем 400/413
	if w.Code == http.StatusBadRequest || w.Code == http.StatusRequestEntityTooLarge {
		t.Errorf("other error should not get body error status, got %d", w.Code)
	}
}

func TestBodyLimitMiddleware_WrapsBody(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mw := BodyLimitMiddleware(10)
	handler := mw(next)
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
}
