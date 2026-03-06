package http

import (
	"net/http"
	"testing"
)

func TestClientIP(t *testing.T) {
	tests := []struct {
		name   string
		xff    string
		xrip   string
		remote string
		wantIP string
	}{
		{"X-Forwarded-For single", "192.168.1.1", "", "1.2.3.4:5678", "192.168.1.1"},
		{"X-Forwarded-For multiple takes first", "10.0.0.1, 10.0.0.2, 10.0.0.3", "", "", "10.0.0.1"},
		{"X-Real-IP when no XFF", "", "172.16.0.1", "", "172.16.0.1"},
		{"RemoteAddr without port", "", "", "192.168.0.1", "192.168.0.1"},
		{"RemoteAddr strips port", "", "", "127.0.0.1:8080", "127.0.0.1"},
		{"XFF trimmed", "  203.0.113.1  ", "", "", "203.0.113.1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &http.Request{
				Header:     http.Header{},
				RemoteAddr: tt.remote,
			}
			if tt.xff != "" {
				r.Header.Set("X-Forwarded-For", tt.xff)
			}
			if tt.xrip != "" {
				r.Header.Set("X-Real-IP", tt.xrip)
			}
			got := ClientIP(r)
			if got != tt.wantIP {
				t.Errorf("ClientIP() = %q, want %q", got, tt.wantIP)
			}
		})
	}
}
