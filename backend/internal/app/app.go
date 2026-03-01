package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"

	"mkk/internal/app/di"
	"mkk/pkg/closer"
	"mkk/pkg/config/contracts"
	healthhttp "mkk/pkg/http/health"
	httprouter "mkk/pkg/http/router"
	httpserver "mkk/pkg/http/server"
	"mkk/pkg/logger"
	"mkk/pkg/metric"
	"mkk/pkg/tracing"
)

// App — приложение: DI, логгер, трейсинг, метрики, HTTP-сервер с роутами.
type App struct {
	cfg        contracts.Provider
	di         *di.Container
	httpServer *http.Server
	listener   net.Listener
	stopRouter func()
}

// New создаёт приложение и инициализирует зависимости (DI, logger, tracing, metrics, DB, Redis, роуты).
// Порядок как в gRPC-примере: DI → инфра → БД → listener → HTTP-сервер и регистрация роутов.
func New(ctx context.Context, cfg contracts.Provider) (*App, error) {
	app := &App{cfg: cfg}

	steps := []func(context.Context) error{
		app.initDI,
		app.initLogger,
		app.initTracer,
		app.initMetrics,
		app.initCloser,
		app.initDatabase,
		app.initListener,
		app.initHTTPServer,
	}
	for _, step := range steps {
		if err := step(ctx); err != nil {
			return nil, err
		}
	}

	return app, nil
}

// Start запускает HTTP-сервер. Блокируется до получения сигнала завершения.
// При корректном shutdown Serve возвращает http.ErrServerClosed — в этом случае ошибку не пробрасываем.
func (a *App) Start(ctx context.Context) error {
	logger.Info(ctx, "🚀 [HTTP] Сервер слушает", zap.String("address", a.listener.Addr().String()))
	err := a.httpServer.Serve(a.listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.di = di.NewContainer(a.cfg)
	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	return logger.Reinit(ctx,
		logger.WithLevel(a.cfg.Logger().Level()),
		logger.WithJSON(a.cfg.Logger().AsJSON()),
		logger.WithName(a.cfg.Logger().Name()),
		logger.WithEnvironment(a.cfg.Logger().Environment()),
		logger.WithOTLPEnable(a.cfg.Logger().OTLPEnable()),
		logger.WithOTLPEndpoint(a.cfg.Logger().OTLPEndpoint()),
		logger.WithOTLPTimeout(a.cfg.Logger().OTLPShutdownTimeout()),
	)
}

func (a *App) initTracer(ctx context.Context) error {
	tracing.SetLogger(logger.Logger())
	if err := tracing.Init(ctx,
		tracing.WithName(a.cfg.App().Name()),
		tracing.WithEnvironment(a.cfg.App().Environment()),
		tracing.WithVersion(a.cfg.App().Version()),
		tracing.WithEnable(a.cfg.Tracing().Enable()),
		tracing.WithEndpoint(a.cfg.Tracing().Endpoint()),
		tracing.WithTimeout(a.cfg.Tracing().Timeout()),
		tracing.WithSampleRatio(a.cfg.Tracing().SampleRatio()),
		tracing.WithRetryEnabled(a.cfg.Tracing().RetryEnabled()),
		tracing.WithRetryInitialInterval(a.cfg.Tracing().RetryInitialInterval()),
		tracing.WithRetryMaxInterval(a.cfg.Tracing().RetryMaxInterval()),
		tracing.WithRetryMaxElapsedTime(a.cfg.Tracing().RetryMaxElapsedTime()),
		tracing.WithTraceContext(a.cfg.Tracing().EnableTraceContext()),
		tracing.WithBaggage(a.cfg.Tracing().EnableBaggage()),
	); err != nil {
		return fmt.Errorf("init tracer: %w", err)
	}
	logger.Info(ctx, "✅ [Tracing] Трейсинг инициализирован", zap.String("service", a.cfg.App().Name()))
	closer.AddNamed("Tracer", func(ctx context.Context) error {
		logger.Info(ctx, "🔍 [Shutdown] Закрытие Tracer")
		return tracing.Shutdown(ctx, a.cfg.Tracing().ShutdownTimeout())
	})
	return nil
}

func (a *App) initMetrics(ctx context.Context) error {
	metric.SetLogger(logger.Logger())
	if err := metric.Init(ctx,
		metric.WithName(a.cfg.App().Name()),
		metric.WithEnvironment(a.cfg.App().Environment()),
		metric.WithVersion(a.cfg.App().Version()),
		metric.WithEnable(a.cfg.Metric().Enable()),
		metric.WithEndpoint(a.cfg.Metric().Endpoint()),
		metric.WithTimeout(a.cfg.Metric().Timeout()),
		metric.WithNamespace(a.cfg.Metric().Namespace()),
		metric.WithAppName(a.cfg.Metric().AppName()),
		metric.WithExportInterval(a.cfg.Metric().ExportInterval()),
		metric.WithShutdownTimeout(a.cfg.Metric().ShutdownTimeout()),
	); err != nil {
		return fmt.Errorf("init metrics: %w", err)
	}
	logger.Info(ctx, "✅ [Metrics] Метрики инициализированы", zap.String("service", a.cfg.App().Name()))
	closer.AddNamed("Metrics", func(ctx context.Context) error {
		logger.Info(ctx, "📊 [Shutdown] Закрытие Metrics")
		return metric.Shutdown(ctx, a.cfg.Metric().ShutdownTimeout())
	})
	return nil
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initDatabase(ctx context.Context) error {
	if _, err := a.di.SqlxDB(ctx); err != nil {
		return fmt.Errorf("mysql pool: %w", err)
	}
	if _, err := a.di.RedisClient(ctx); err != nil {
		return fmt.Errorf("redis client: %w", err)
	}
	return nil
}

func (a *App) initListener(_ context.Context) error {
	addr := a.cfg.HTTP().Address()
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", addr, err)
	}
	a.listener = listener

	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		logger.Info(ctx, "🔌 [Shutdown] Закрытие listener")
		if err := listener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			return err
		}
		return nil
	})
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	addr := a.cfg.HTTP().Address()
	timeout := a.cfg.HTTP().Timeout()
	buckets := a.cfg.Metric().BucketBoundaries()
	if len(buckets) == 0 {
		buckets = nil
	}

	r, stopRouter := httprouter.NewRouter(ctx, timeout, nil, nil, nil, nil, false, 0, buckets)
	a.stopRouter = stopRouter

	healthhttp.RegisterRoutes(r)

	if err := a.di.RegisterAccountRoutes(ctx, r); err != nil {
		return fmt.Errorf("register account routes: %w", err)
	}

	a.httpServer = httpserver.NewServer(r, addr, 5*time.Second, timeout, timeout, 60*time.Second, 1<<20)

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		logger.Info(ctx, "⚡ [Shutdown] Остановка HTTP сервера")
		a.stopRouter()
		return a.httpServer.Shutdown(ctx)
	})

	return nil
}
