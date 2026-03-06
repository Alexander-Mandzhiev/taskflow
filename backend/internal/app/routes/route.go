package routes

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	account_v1 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/account/v1"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/app/routes/public"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/app/routes/session_auth"
	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http/middleware"
)

// Middlewares содержит middleware для роутов (JWT auth, body limit, user rate limit).
// StopUserRateLimit вызывать при graceful shutdown (main).
type Middlewares struct {
	BodyLimit         func(http.Handler) http.Handler
	JWT               *middleware.JWTAuthMiddleware
	UserRateLimit     func(http.Handler) http.Handler
	StopUserRateLimit func()
}

// NewMiddlewares создаёт middleware для публичных и защищённых роутов (JWT).
func NewMiddlewares(
	ctx context.Context,
	accessSecret string,
	isSecure bool,
	cookieDomain string,
	accessTokenCookieName string,
) *Middlewares {
	userRateLimitMw, stopUserLimiter := middleware.UserRateLimitMiddleware(ctx)

	return &Middlewares{
		BodyLimit:         middleware.BodyLimitMiddleware(pkghttp.MaxRequestBodyBytes),
		JWT:               middleware.NewJWTAuthMiddleware(accessSecret, isSecure, cookieDomain, accessTokenCookieName),
		UserRateLimit:     userRateLimitMw,
		StopUserRateLimit: stopUserLimiter,
	}
}

// RegisterAPIs регистрирует все API account под префиксом /api/v1:
// публичные (register, login); защищённые (logout, whoami) с JWT и UserRateLimit.
func RegisterAPIs(ctx context.Context, router *chi.Mux, api *account_v1.API, mw *Middlewares) {
	router.Route("/api/v1", func(r chi.Router) {
		r.Use(mw.BodyLimit)

		r.Group(func(r chi.Router) {
			public.Register(ctx, r, api)
		})

		r.Group(func(r chi.Router) {
			r.Use(mw.JWT.Handle)
			r.Use(mw.UserRateLimit)
			session_auth.Register(ctx, r, api)
		})
	})
}
