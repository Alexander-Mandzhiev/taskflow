package di

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	account_v1 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/account/v1"
	task_v1 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/task/v1"
	team_v1 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/team/v1"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/app/routes/public"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/app/routes/session_auth"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/app/routes/task"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/app/routes/team"
	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http/middleware"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

// accountMiddlewares — набор middleware для роутов account. Создаётся лениво в getAccountMiddlewares.
type accountMiddlewares struct {
	bodyLimit         func(http.Handler) http.Handler
	jwt               *middleware.JWTAuthMiddleware
	userRateLimit     func(http.Handler) http.Handler
	stopUserRateLimit func()
}

// RegisterAccountRoutes регистрирует API account, team и task на роутере. Middleware создаются лениво при первом вызове.
func (d *Container) RegisterAccountRoutes(ctx context.Context, router *chi.Mux) error {
	if err := d.requireCloser(); err != nil {
		return err
	}
	api, err := d.AccountV1API(ctx)
	if err != nil {
		return err
	}
	teamAPI, err := d.TeamV1API(ctx)
	if err != nil {
		return err
	}
	taskAPI, err := d.TaskV1API(ctx)
	if err != nil {
		return err
	}
	mw, err := d.getAccountMiddlewares(ctx)
	if err != nil {
		return err
	}
	registerAccountRoutes(router, ctx, api, teamAPI, taskAPI, mw)
	return nil
}

// getAccountMiddlewares возвращает middleware для account (ленивая инициализация, кеш в контейнере).
func (d *Container) getAccountMiddlewares(ctx context.Context) (*accountMiddlewares, error) {
	if d.accountMiddlewares != nil {
		return d.accountMiddlewares, nil
	}
	accountSvc, err := d.AccountService(ctx)
	if err != nil {
		return nil, fmt.Errorf("account service: %w", err)
	}
	jwtCfg := d.cfg.JWT()
	sessionCfg := d.cfg.Session()

	jwtMw := middleware.NewJWTAuthMiddleware(
		jwtCfg.AccessSecret(),
		sessionCfg.IsSecure(),
		sessionCfg.CookieDomain(),
		jwtCfg.AccessTokenCookieName(),
		accountSvc,
		jwtCfg.RefreshTokenCookieName(),
		jwtCfg.AccessTTL(),
	)
	userRateLimitMw, stopUserRateLimit := middleware.UserRateLimitMiddleware(ctx)

	d.accountMiddlewares = &accountMiddlewares{
		bodyLimit:         middleware.BodyLimitMiddleware(pkghttp.MaxRequestBodyBytes),
		jwt:               jwtMw,
		userRateLimit:     userRateLimitMw,
		stopUserRateLimit: stopUserRateLimit,
	}
	d.cl.Add(func(ctx context.Context) error {
		d.accountMiddlewares.stopUserRateLimit()
		logger.Info(ctx, "🚦 [Shutdown] Closed User rate limiter")
		return nil
	})
	return d.accountMiddlewares, nil
}

// registerAccountRoutes вешает middleware и регистрирует группы путей на роутере.
func registerAccountRoutes(router *chi.Mux, ctx context.Context, api *account_v1.API, teamAPI *team_v1.API, taskAPI *task_v1.API, mw *accountMiddlewares) {
	router.Route("/api/v1", func(r chi.Router) {
		r.Use(mw.bodyLimit)
		r.Group(func(r chi.Router) {
			registerAccountPublicGroup(r, ctx, api)
		})
		r.Group(func(r chi.Router) {
			r.Use(mw.jwt.Handle)
			// TODO: для нагрузочного тестирования user rate limit отключён; перед продом вернуть: r.Use(mw.userRateLimit)
			registerAccountPrivateGroup(r, ctx, api, teamAPI, taskAPI)
		})
	})
}

// registerAccountPublicGroup регистрирует публичные роуты (register, login).
func registerAccountPublicGroup(r chi.Router, ctx context.Context, api *account_v1.API) {
	public.Register(ctx, r, api)
}

// registerAccountPrivateGroup регистрирует защищённые JWT роуты (logout, teams, tasks, reports).
func registerAccountPrivateGroup(r chi.Router, ctx context.Context, api *account_v1.API, teamAPI *team_v1.API, taskAPI *task_v1.API) {
	session_auth.Register(ctx, r, api)
	team.Register(ctx, r, teamAPI)
	task.Register(ctx, r, taskAPI)
}
