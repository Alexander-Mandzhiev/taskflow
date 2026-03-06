package tracing

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

// Tracing управляет состоянием трассировки и её инициализацией
type Tracing struct {
	initOnce       sync.Once
	tracerProvider *sdktrace.TracerProvider
	logger         Logger
}

// Глобальный экземпляр для использования по всему приложению
var globalTracing = NewWithLogger(&logger.NoopLogger{})

// SetLogger устанавливает логгер для глобального экземпляра трассировки
func SetLogger(logger Logger) {
	globalTracing.SetLogger(logger)
}

// NewWithLogger создает новый экземпляр Tracing с указанным логгером.
// Если logger == nil, используется NoopLogger.
func NewWithLogger(l Logger) *Tracing {
	if l == nil {
		l = &logger.NoopLogger{}
	}
	return &Tracing{
		logger: l,
	}
}

// SetLogger устанавливает логгер для Tracing. Если l == nil, не меняет текущий логгер.
func (t *Tracing) SetLogger(l Logger) {
	if l != nil {
		t.logger = l
	}
}

const (
	// DefaultCompressor - алгоритм сжатия по умолчанию
	DefaultCompressor = "gzip"
)

// Init инициализирует OpenTelemetry TracerProvider и все инструменты трассировки
func Init(ctx context.Context, opts ...Option) error {
	return globalTracing.Init(ctx, opts...)
}

// Init инициализирует OpenTelemetry TracerProvider и все инструменты трассировки для конкретного экземпляра
func (t *Tracing) Init(ctx context.Context, opts ...Option) error {
	var initErr error

	t.initOnce.Do(func() {
		cfg := defaultConfig()
		for _, opt := range opts {
			opt(cfg)
		}

		if !cfg.enable {
			return // Трейсинг отключен - tracerProvider остается nil
		}

		initErr = t.initTracer(ctx, cfg)
		if initErr != nil {
			return
		}
	})

	return initErr
}

// Reinit сбрасывает состояние и инициализирует трейсер заново с функциональными опциями.
func Reinit(ctx context.Context, opts ...Option) error {
	return globalTracing.Reinit(ctx, opts...)
}

// Reinit сбрасывает состояние и инициализирует трейсер заново для конкретного экземпляра
func (t *Tracing) Reinit(ctx context.Context, opts ...Option) error {
	t.initOnce = sync.Once{}
	return t.Init(ctx, opts...)
}

// ParseDuration парсит строку в time.Duration с дефолтным значением.
func ParseDuration(s string, defaultValue time.Duration) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return defaultValue
	}
	return d
}

// initTracer выполняет фактическую инициализацию трассировки для конкретного экземпляра
func (t *Tracing) initTracer(ctx context.Context, cfg *Config) error {
	// 1. Создаем OTLP gRPC экспортер для отправки трейсов в OpenTelemetry Collector
	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(cfg.endpoint), // Адрес коллектора
		otlptracegrpc.WithInsecure(),             // Отключаем TLS для локальной разработки
		otlptracegrpc.WithTimeout(cfg.timeout),
		otlptracegrpc.WithCompressor(DefaultCompressor),
		otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{
			Enabled:         cfg.retryEnabled,
			InitialInterval: cfg.retryInitialInterval,
			MaxInterval:     cfg.retryMaxInterval,
			MaxElapsedTime:  cfg.retryMaxElapsedTime,
		}),
	)
	if err != nil {
		t.logger.Error(ctx, "[Tracing] Failed to create OTLP exporter", zap.Error(err))
		return err
	}

	// 2. Создаем ресурс с метаданными о сервисе
	attributeResource, err := resource.New(ctx,
		resource.WithAttributes(
			// Используем стандартные атрибуты OpenTelemetry
			semconv.ServiceName(cfg.name),
			semconv.ServiceVersion(cfg.version),
			attribute.String("environment", cfg.environment),
		),
		// Автоматически определяем хост, ОС и другие системные атрибуты
		resource.WithHost(),
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithContainer(),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		t.logger.Error(ctx, "[Tracing] Failed to create resource", zap.Error(err))
		return err
	}

	// 3. Создаем TracerProvider с настроенным экспортером и ресурсом
	t.tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(attributeResource),
		// Настраиваем семплирование трейсов:
		// 1. ParentBased - учитываем решение о семплировании родительского спана
		// 2. TraceIDRatioBased - используем настройку из конфига
		// В продакшене рекомендуется использовать меньший процент (0.1 = 10%)
		// для снижения нагрузки на систему трассировки
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.sampleRatio/100.0))),
	)

	// Устанавливаем глобальный провайдер трейсов
	otel.SetTracerProvider(t.tracerProvider)

	// Настраиваем пропагацию контекста для передачи между сервисами:
	// 1. TraceContext - стандарт W3C для передачи trace ID и parent span ID через HTTP заголовки
	//    Позволяет связать запросы между сервисами в единый трейс
	// 2. Baggage - механизм для передачи дополнительных метаданных между сервисами
	//    Например: user_id, tenant_id, request_id и другие бизнес-контексты
	// Пропагация - это механизм передачи контекста трассировки между сервисами
	// Когда запрос проходит через несколько сервисов, пропагация позволяет:
	// - Сохранить связь между всеми спанами в цепочке вызовов
	// - Передавать дополнительный контекст между сервисами
	// - Обеспечить сквозную трассировку всего запроса

	var propagators []propagation.TextMapPropagator

	if cfg.enableTraceContext {
		propagators = append(propagators, propagation.TraceContext{})
	}

	if cfg.enableBaggage {
		propagators = append(propagators, propagation.Baggage{})
	}

	if len(propagators) == 0 {
		propagators = append(propagators, propagation.TraceContext{})
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagators...))

	return nil
}

// GetTracerProvider возвращает текущий провайдер трейсов
func GetTracerProvider() *sdktrace.TracerProvider {
	return globalTracing.GetTracerProvider()
}

// GetTracerProvider возвращает текущий провайдер трейсов для конкретного экземпляра
func (t *Tracing) GetTracerProvider() *sdktrace.TracerProvider {
	if t.tracerProvider == nil {
		return sdktrace.NewTracerProvider()
	}
	return t.tracerProvider
}

// Shutdown закрывает провайдер трейсов в правильном порядке.
func Shutdown(ctx context.Context, timeout time.Duration) error {
	return globalTracing.Shutdown(ctx, timeout)
}

// Shutdown закрывает провайдер трейсов в правильном порядке для конкретного экземпляра.
// Используем независимый контекст (не ctx), чтобы гарантировать полный timeout на shutdown.
//
//nolint:contextcheck // shutdown намеренно имеет свой таймаут, не наследуем ctx
func (t *Tracing) Shutdown(ctx context.Context, timeout time.Duration) error {
	if t.tracerProvider == nil {
		return nil
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Закрываем провайдер трейсов при выходе
	err := t.tracerProvider.Shutdown(shutdownCtx)
	if err != nil {
		return err
	}

	t.tracerProvider = nil

	return nil
}
