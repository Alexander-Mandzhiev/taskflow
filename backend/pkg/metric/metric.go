package metric

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.uber.org/zap"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

// Metrics управляет состоянием метрик и их инициализацией
type Metrics struct {
	initOnce      sync.Once
	exporter      *otlpmetricgrpc.Exporter
	meterProvider *sdkmetric.MeterProvider
	logger        Logger
	config        *Config
}

// Глобальный экземпляр для использования по всему приложению
var globalMetrics = NewWithLogger(&logger.NoopLogger{})

// SetLogger устанавливает логгер для глобального экземпляра метрик
func SetLogger(logger Logger) {
	globalMetrics.SetLogger(logger)
}

// getMetricName генерирует полное имя метрики с namespace и appName для конкретного экземпляра
func (m *Metrics) getMetricName(metricName string) string {
	if m.config == nil {
		return metricName
	}
	return m.config.namespace + "_" + m.config.appName + "_" + metricName
}

// NewWithLogger создает новый экземпляр Metrics с указанным логгером
func NewWithLogger(l Logger) *Metrics {
	if l == nil {
		l = &logger.NoopLogger{}
	}
	return &Metrics{
		logger: l,
	}
}

// SetLogger устанавливает логгер для Metrics
func (m *Metrics) SetLogger(logger Logger) {
	if logger != nil {
		m.logger = logger
	}
}

// Init инициализирует OpenTelemetry MeterProvider и все инструменты метрик
func Init(ctx context.Context, opts ...Option) error {
	return globalMetrics.Init(ctx, opts...)
}

// Init инициализирует OpenTelemetry MeterProvider и все инструменты метрик для конкретного экземпляра
// Важно: может быть вызван только один раз для экземпляра Metrics. Повторные вызовы будут проигнорированы.
func (m *Metrics) Init(ctx context.Context, opts ...Option) error {
	var initErr error

	m.initOnce.Do(func() {
		// Защита от nil logger
		if m.logger == nil {
			m.logger = &logger.NoopLogger{}
		}

		cfg := defaultConfig()
		for _, opt := range opts {
			opt(cfg)
		}

		// Сохраняем конфигурацию в структуре для доступа к namespace и appName
		m.config = cfg

		if !cfg.enable {
			m.logInfo(ctx, "Metrics disabled")
			return // Метрики отключены - meterProvider остается nil
		}

		initErr = m.initMetrics(ctx, cfg)
		if initErr != nil {
			return
		}
	})

	return initErr
}

// initMetrics выполняет фактическую инициализацию метрик для конкретного экземпляра
func (m *Metrics) initMetrics(ctx context.Context, cfg *Config) error {
	var err error

	m.exporter, err = otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(cfg.endpoint),
		otlpmetricgrpc.WithTLSCredentials(insecure.NewCredentials()),
		otlpmetricgrpc.WithTimeout(cfg.timeout),
	)
	if err != nil {
		m.logError(ctx, "failed to create OTLP exporter", err)
		return errors.Wrap(err, "failed to create OTLP exporter")
	}

	// 2. Создаем ресурс с метаданными о сервисе
	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("service.name", cfg.name),
			attribute.String("service.version", cfg.version),
			attribute.String("deployment.environment", cfg.environment),
		),
	)
	if err != nil {
		m.logError(ctx, "failed to create resource", err)
		return errors.Wrap(err, "failed to create resource")
	}

	// 3. Создаем MeterProvider с периодическим reader
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				m.exporter,
				sdkmetric.WithInterval(cfg.exportInterval),
			),
		),
	)

	m.meterProvider = meterProvider
	otel.SetMeterProvider(meterProvider)

	m.logInfo(ctx, "Metrics initialized successfully")
	return nil
}

// GetMeterProvider возвращает текущий провайдер метрик
func GetMeterProvider() *sdkmetric.MeterProvider {
	return globalMetrics.GetMeterProvider()
}

// GetMeterProvider возвращает текущий провайдер метрик для конкретного экземпляра.
// Возвращает nil, если метрики не инициализированы.
func (m *Metrics) GetMeterProvider() *sdkmetric.MeterProvider {
	return m.meterProvider
}

// Shutdown закрывает провайдер метрик и экспортер в правильном порядке.
// MeterProvider должен закрываться первым, чтобы корректно завершить отправку данных в экспортер.
func Shutdown(ctx context.Context, timeout time.Duration) error {
	return globalMetrics.Shutdown(ctx, timeout)
}

// Shutdown закрывает провайдер метрик и экспортер в правильном порядке для конкретного экземпляра.
// Используем независимый контекст (не ctx), чтобы гарантировать полный timeout на shutdown,
// даже если переданный ctx уже отменён (например при SIGTERM).
//
//nolint:contextcheck // shutdown намеренно имеет свой таймаут, не наследуем ctx
func (m *Metrics) Shutdown(ctx context.Context, timeout time.Duration) error {
	if m.meterProvider == nil && m.exporter == nil {
		return nil
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var err error

	// 1. Сначала закрываем MeterProvider - он может продолжать отправлять данные в экспортер
	if m.meterProvider != nil {
		err = m.meterProvider.Shutdown(shutdownCtx)
		if err != nil {
			m.logError(shutdownCtx, "failed to shutdown meter provider", err)
			return errors.Wrap(err, "failed to shutdown meter provider")
		}
		m.meterProvider = nil
	}

	// 2. Затем закрываем экспортер - после того как MeterProvider завершил отправку данных
	// MeterProvider может автоматически закрыть exporter через reader, поэтому проверяем
	if m.exporter != nil {
		err = m.exporter.Shutdown(shutdownCtx)
		if err != nil {
			if isExporterAlreadyShutdownError(err) {
				m.logInfo(shutdownCtx, "Exporter already shutdown, ignoring")
			} else {
				m.logError(shutdownCtx, "failed to shutdown exporter", err)
				return errors.Wrap(err, "failed to shutdown exporter")
			}
		}
		m.exporter = nil
	}

	return nil
}

// isExporterAlreadyShutdownError проверяет, является ли ошибка признаком того, что exporter уже закрыт
func isExporterAlreadyShutdownError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, "gRPC exporter is shutdown") ||
		strings.Contains(errMsg, "exporter is already shutdown")
}

// logError логирует ошибку, если logger инициализирован
func (m *Metrics) logError(ctx context.Context, msg string, err error) {
	if m.logger != nil {
		m.logger.Error(ctx, msg, zap.Error(err))
	}
}

// logInfo логирует информационное сообщение, если logger инициализирован
func (m *Metrics) logInfo(ctx context.Context, msg string) {
	if m.logger != nil {
		m.logger.Info(ctx, msg)
	}
}

// IsInitialized возвращает true, если метрики инициализированы
func (m *Metrics) IsInitialized() bool {
	return m.meterProvider != nil
}

// IsEnabled возвращает true, если метрики включены в конфигурации
func (m *Metrics) IsEnabled() bool {
	return m.config != nil && m.config.enable
}

// GetConfig возвращает текущую конфигурацию метрик (может быть nil)
func (m *Metrics) GetConfig() *Config {
	return m.config
}

// HealthCheck проверяет состояние метрик и возвращает ошибку, если что-то не так
func (m *Metrics) HealthCheck(ctx context.Context) error {
	if m.meterProvider == nil {
		return errors.New("metrics not initialized")
	}

	if m.exporter == nil {
		return errors.New("metrics exporter not initialized")
	}

	return nil
}
