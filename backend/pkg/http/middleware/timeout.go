package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"
)

// TimeoutByPath возвращает middleware с таймаутом: для путей с префиксом debugPrefix — debugTimeout, иначе normalTimeout.
// Нужно для /debug/pprof/profile?seconds=30: сбор 30s не должен обрываться по normalTimeout (30s).
func TimeoutByPath(normalTimeout, debugTimeout time.Duration, debugPrefix string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := normalTimeout
			if debugPrefix != "" && strings.HasPrefix(r.URL.Path, debugPrefix) {
				t = debugTimeout
			}
			ctx, cancel := context.WithTimeout(r.Context(), t)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
