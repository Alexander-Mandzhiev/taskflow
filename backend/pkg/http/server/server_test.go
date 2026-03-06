package http_server

import (
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

func TestNewServer_ReturnsNonNil(t *testing.T) {
	r := chi.NewRouter()
	srv := NewServer(r, ":8080", 0, 0, 0, 0, 0)
	if srv == nil {
		t.Fatal("NewServer returned nil")
	}
	if srv.Handler != r {
		t.Error("Handler should be the passed router")
	}
}

func TestNewServer_SetsFields(t *testing.T) {
	r := chi.NewRouter()
	addr := ":9090"
	readHeader := 2 * time.Second
	read := 10 * time.Second
	write := 15 * time.Second
	idle := 30 * time.Second
	maxHeader := 4096

	srv := NewServer(r, addr, readHeader, read, write, idle, maxHeader)

	if srv.Addr != addr {
		t.Errorf("Addr = %q, want %q", srv.Addr, addr)
	}
	if srv.Handler != r {
		t.Error("Handler != router")
	}
	if srv.ReadHeaderTimeout != readHeader {
		t.Errorf("ReadHeaderTimeout = %v, want %v", srv.ReadHeaderTimeout, readHeader)
	}
	if srv.ReadTimeout != read {
		t.Errorf("ReadTimeout = %v, want %v", srv.ReadTimeout, read)
	}
	if srv.WriteTimeout != write {
		t.Errorf("WriteTimeout = %v, want %v", srv.WriteTimeout, write)
	}
	if srv.IdleTimeout != idle {
		t.Errorf("IdleTimeout = %v, want %v", srv.IdleTimeout, idle)
	}
	if srv.MaxHeaderBytes != maxHeader {
		t.Errorf("MaxHeaderBytes = %d, want %d", srv.MaxHeaderBytes, maxHeader)
	}
}

func TestNewServer_ZeroValues(t *testing.T) {
	r := chi.NewRouter()
	srv := NewServer(r, "", 0, 0, 0, 0, 0)

	if srv == nil {
		t.Fatal("NewServer returned nil")
	}
	if srv.Handler != r {
		t.Error("Handler != router")
	}
	// Сервер с нулевыми таймаутами допустим (значения по умолчанию http.Server)
	if srv.Addr != "" {
		t.Errorf("Addr = %q, want empty", srv.Addr)
	}
}

func TestNewServer_HandlerIsChiMux(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		data := []byte("pong")
		n, err := w.Write(data)
		if err != nil {
			t.Errorf("w.Write: %v", err)
		}
		if n != len(data) {
			t.Errorf("w.Write wrote %d bytes, want %d", n, len(data))
		}
	})

	srv := NewServer(r, "127.0.0.1:0", 0, time.Second, time.Second, 0, 0)
	if srv.Handler == nil {
		t.Fatal("Handler is nil")
	}
	_, ok := srv.Handler.(*chi.Mux)
	if !ok {
		t.Error("Handler should be *chi.Mux")
	}
}
