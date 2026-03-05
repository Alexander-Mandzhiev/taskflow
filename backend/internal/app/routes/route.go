package routes

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	account_v1 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/account/v1"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/app/routes/public"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/app/routes/session_auth"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/service"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http/middleware"
)

// maxRequestBodyBytes — лимит тела запроса для POST, защита от исчерпания памяти.
const maxRequestBodyBytes = 1 << 20 // 1 MiB

// Middlewares содержит middleware для роутов (auth, body limit, session, user rate limit).
// StopUserRateLimit вызывать при graceful shutdown (main).
type Middlewares struct {
	BodyLimit         func(http.Handler) http.Handler
	Session           func(http.Handler) http.Handler
	Auth              *middleware.AuthMiddleware
	UserRateLimit     func(http.Handler) http.Handler
	StopUserRateLimit func()
}

// NewMiddlewares создаёт middleware для публичных и защищённых роутов.
// sessionService используется в AuthMiddleware для проверки сессии (Whoami).
func NewMiddlewares(
	sessionService service.AccountService,
	isSecure bool,
	cookieDomain string,
) *Middlewares {
	userRateLimitMw, stopUserLimiter := middleware.UserRateLimitMiddleware()
	return &Middlewares{
		BodyLimit:         middleware.BodyLimitMiddleware(maxRequestBodyBytes),
		Session:           middleware.SessionMiddleware,
		Auth:              middleware.NewAuthMiddleware(sessionService, isSecure, cookieDomain),
		UserRateLimit:     userRateLimitMw,
		StopUserRateLimit: stopUserLimiter,
	}
}

// RegisterAPIs регистрирует все API account под префиксом /api/v1:
// публичные роуты (register, login) без auth; защищённые (logout) с Auth и UserRateLimit.
func RegisterAPIs(ctx context.Context, router *chi.Mux, api *account_v1.API, mw *Middlewares) {
	router.Route("/api/v1", func(r chi.Router) {
		r.Use(mw.Session)
		r.Use(mw.BodyLimit) // на весь /api/v1 — защита и публичных, и приватных POST

		// 1. Публичные роуты (без auth)
		r.Group(func(r chi.Router) {
			public.Register(ctx, r, api)
		})

		// 2. Роуты с session auth (Auth + лимит на авторизованных пользователей)
		r.Group(func(r chi.Router) {
			r.Use(mw.Auth.Handle)
			r.Use(mw.UserRateLimit)
			session_auth.Register(ctx, r, api)
		})
	})
}
