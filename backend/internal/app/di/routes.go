package di

import (
	"context"

	"github.com/go-chi/chi/v5"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/app/routes"
)

// RegisterAccountRoutes регистрирует API account (register, login, logout, whoami) на роутере и добавляет
// остановку user rate limiter в closer для graceful shutdown.
func (d *Container) RegisterAccountRoutes(ctx context.Context, router *chi.Mux) error {
	if err := d.requireCloser(); err != nil {
		return err
	}
	api, err := d.AccountV1API(ctx)
	if err != nil {
		return err
	}
	accountSvc, err := d.AccountService(ctx)
	if err != nil {
		return err
	}
	sessionCfg := d.cfg.Session()
	mw := routes.NewMiddlewares(accountSvc, sessionCfg.IsSecure(), sessionCfg.CookieDomain())
	routes.RegisterAPIs(ctx, router, api, mw)

	d.cl.AddNamed("User rate limiter", func(ctx context.Context) error {
		mw.StopUserRateLimit()
		return nil
	})
	return nil
}
