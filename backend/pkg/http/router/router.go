// Package http_router предоставляет функции для создания HTTP роутера.
//
// Роутер навешивает chi (RequestID, RealIP, Recoverer, Timeout) и переданный
// слайс глобальных middleware. Сборка глобальных middleware — в app.initHTTPRouter.
package http_router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http/middleware"
)

// NewRouter создаёт роутер: chi (RequestID, RealIP) → globalMiddlewares → Recoverer, Timeout.
// globalMiddlewares задаются снаружи (в app), порядок в слайсе сохраняется.
// Один timeout из конфига для всех путей (API и /debug/); в local — больше для pprof, в проде — меньше.
func NewRouter(timeout time.Duration, globalMiddlewares []func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)

	for _, mw := range globalMiddlewares {
		r.Use(mw)
	}

	r.Use(chimw.Recoverer)
	r.Use(middleware.TimeoutByPath(timeout, timeout, "/debug/"))

	return r
}
