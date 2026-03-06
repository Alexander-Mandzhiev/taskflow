package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/app/di"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/closer"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/contracts"
	healthhttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http/health"
	httprouter "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http/router"
	httpserver "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http/server"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metric"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/tracing"
)

// App — приложение: DI, логгер, трейсинг, метрики, HTTP-сервер с роутами.
// Каждый экземпляр владеет своим closer — тесты можно запускать параллельно и изолированно.
type App struct {
	cfg        contracts.Provider
	di         *di.Container
	closer     *closer.Closer
	httpServer *http.Server
	listener   net.Listener
	stopRouter func()
}

// New создаёт приложение и инициализирует зависимости (DI, logger, tracing, metrics, DB, Redis, роуты).
// cl — менеджер ресурсов для graceful shutdown; создаётся снаружи (main/тест), чтобы точка входа владела жизненным циклом.
// Порядок как в gRPC-примере: DI → инфра → БД → listener → HTTP-сервер и регистрация роутов.
func New(ctx context.Context, cfg contracts.Provider, cl *closer.Closer) (*App, error) {
	if cl == nil {
		return nil, fmt.Errorf("closer must not be nil")
	}
	app := &App{cfg: cfg, closer: cl}

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

// Shutdown закрывает все зарегистрированные ресурсы (БД, Redis, HTTP, tracer, metrics и т.д.).
// Вызывать при завершении приложения (например из defer в main). Идемпотентен.
func (a *App) Shutdown(ctx context.Context) error {
	return a.closer.CloseAll(ctx)
}

func (a *App) initDI(_ context.Context) error {
	a.di = di.NewContainer(a.cfg)
	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	if err := logger.Reinit(ctx,
		logger.WithLevel(a.cfg.Logger().Level()),
		logger.WithJSON(a.cfg.Logger().AsJSON()),
		logger.WithName(a.cfg.Logger().Name()),
		logger.WithEnvironment(a.cfg.Logger().Environment()),
		logger.WithOTLPEnable(a.cfg.Logger().OTLPEnable()),
		logger.WithOTLPEndpoint(a.cfg.Logger().OTLPEndpoint()),
		logger.WithOTLPShutdownTimeout(a.cfg.Logger().OTLPShutdownTimeout()),
	); err != nil {
		return err
	}

	a.closer.Add(func(ctx context.Context) error {
		err := logger.Shutdown(ctx, a.cfg.Logger().OTLPShutdownTimeout())
		logger.Info(ctx, "📝 [Shutdown] Closed Logger")
		return err
	})

	logger.Info(ctx, "✅ [Logger] Логгер инициализирован", zap.String("name", a.cfg.Logger().Name()), zap.String("level", a.cfg.Logger().Level()))
	return nil
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
	a.closer.Add(func(ctx context.Context) error {
		err := tracing.Shutdown(ctx, a.cfg.Tracing().ShutdownTimeout())
		logger.Info(ctx, "🔍 [Shutdown] Closed Tracer")
		return err
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
	a.closer.Add(func(ctx context.Context) error {
		err := metric.Shutdown(ctx, a.cfg.Metric().ShutdownTimeout())
		logger.Info(ctx, "📊 [Shutdown] Closed Metrics")
		return err
	})
	return nil
}

func (a *App) initCloser(ctx context.Context) error {
	a.di.SetCloser(a.closer)
	logger.Info(ctx, "✅ [Closer] Менеджер graceful shutdown подключён к App")
	return nil
}

func (a *App) initDatabase(ctx context.Context) error {
	if _, err := a.di.SqlxDB(ctx); err != nil {
		return fmt.Errorf("mysql pool: %w", err)
	}
	logger.Info(ctx, "✅ [Database] MySQL пул создан и проверен")

	if err := a.di.RunMigrations(ctx); err != nil {
		return fmt.Errorf("migrations: %w", err)
	}

	if _, err := a.di.RedisClient(ctx); err != nil {
		return fmt.Errorf("redis client: %w", err)
	}
	logger.Info(ctx, "✅ [Database] Redis клиент создан и проверен")
	return nil
}

func (a *App) initListener(ctx context.Context) error {
	addr := a.cfg.HTTP().Address()
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", addr, err)
	}
	a.listener = listener
	logger.Info(ctx, "✅ [HTTP] Listener создан", zap.String("address", addr))

	a.closer.Add(func(ctx context.Context) error {
		err := listener.Close()
		if err != nil && !errors.Is(err, net.ErrClosed) {
			return err
		}
		logger.Info(ctx, "🔌 [Shutdown] Closed TCP listener")
		return nil
	})
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	httpCfg := a.cfg.HTTP()
	addr := httpCfg.Address()
	timeout := httpCfg.Timeout()
	buckets := a.cfg.Metric().BucketBoundaries()
	if len(buckets) == 0 {
		buckets = nil
	}

	corsCfg := a.cfg.CORS()
	r, stopRouter := httprouter.NewRouter(ctx, timeout,
		corsCfg.AllowedOrigins(),
		corsCfg.AllowedMethods(),
		corsCfg.AllowedHeaders(),
		corsCfg.ExposedHeaders(),
		corsCfg.AllowCredentials(),
		corsCfg.MaxAge(),
		buckets,
	)
	a.stopRouter = stopRouter

	healthhttp.RegisterRoutes(r)

	if err := a.di.RegisterAccountRoutes(ctx, r); err != nil {
		return fmt.Errorf("register account routes: %w", err)
	}

	a.httpServer = httpserver.NewServer(r, addr,
		httpCfg.ReadHeaderTimeout(),
		httpCfg.ReadTimeout(),
		httpCfg.WriteTimeout(),
		httpCfg.IdleTimeout(),
		httpCfg.MaxHeaderBytes(),
	)

	a.closer.Add(func(ctx context.Context) error {
		a.stopRouter()
		shutdownCtx, cancel := context.WithTimeout(ctx, httpCfg.ShutdownTimeout())
		defer cancel()
		err := a.httpServer.Shutdown(shutdownCtx)
		logger.Info(ctx, "⚡ [Shutdown] Closed HTTP server")
		return err
	})

	logger.Info(ctx, "✅ [HTTP] Роуты зарегистрированы, сервер готов к запуску", zap.String("address", addr))
	return nil
}
