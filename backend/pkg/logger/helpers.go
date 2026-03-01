package logger

import (
	"context"
	"os"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	otelLog "go.opentelemetry.io/otel/log"
	otelLogSdk "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Ключи для извлечения данных из контекста
type contextKey string

const (
	traceIDKey   contextKey = "trace_id"   // Глобальный идентификатор трассировки запроса
	requestIDKey contextKey = "request_id" // Уникальный идентификатор HTTP/gRPC запроса
	userIDKey    contextKey = "user_id"    // Идентификатор пользователя (для будущего использования)
)

// buildCores создает слайс cores для zapcore.Tee.
// Всегда включает stdout core, опционально добавляет OTLP core.
func buildCores(ctx context.Context, cfg *Config) []zapcore.Core {
	cores := []zapcore.Core{
		createStdoutCore(cfg.asJSON),
	}

	if cfg.otlp.enable {
		if otlpCore := createOTLPCore(ctx, cfg); otlpCore != nil {
			cores = append(cores, otlpCore)
		}
	}

	return cores
}

// createStdoutCore создает core для записи в stdout/stderr.
// Поддерживает JSON и консольный формат вывода.
func createStdoutCore(asJSON bool) zapcore.Core {
	config := buildEncoderConfig()
	var encoder zapcore.Encoder
	if asJSON {
		encoder = zapcore.NewJSONEncoder(config)
	} else {
		encoder = zapcore.NewConsoleEncoder(config)
	}

	return zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
}

// buildEncoderConfig настраивает формат вывода логов с нужными полями
func buildEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:      "timestamp",
		LevelKey:     "level",
		MessageKey:   "message",
		CallerKey:    "caller",
		LineEnding:   zapcore.DefaultLineEnding,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
}

// createOTLPCore создает core для отправки в OpenTelemetry коллектор.
// При ошибке подключения возвращает nil (graceful degradation).
func createOTLPCore(ctx context.Context, cfg *Config) *SimpleOTLPCore {
	otlpLogger, err := createOTLPLogger(ctx, cfg)
	if err != nil {
		// Логирование ошибки невозможно, так как логгер еще не инициализирован
		return nil
	}

	// Прямо передаём OTLP-логгер в core. Буферизацию делает OTLP SDK (BatchProcessor).
	return NewSimpleOTLPCore(otlpLogger, level)
}

// createOTLPLogger создает OTLP логгер с настроенным экспортером и ресурсами.
// Использует BatchProcessor для эффективной отправки логов.
func createOTLPLogger(ctx context.Context, cfg *Config) (otelLog.Logger, error) {
	exporter, err := createOTLPExporter(ctx, cfg.otlp.endpoint)
	if err != nil {
		return nil, err
	}

	rs, err := createResource(ctx, cfg.name, cfg.environment)
	if err != nil {
		return nil, err
	}

	provider := otelLogSdk.NewLoggerProvider(
		otelLogSdk.WithResource(rs),
		otelLogSdk.WithProcessor(otelLogSdk.NewBatchProcessor(exporter)),
	)
	otelProvider = provider

	return provider.Logger("app"), nil
}

// createResource создает метаданные сервиса для телеметрии
func createResource(ctx context.Context, name, environment string) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(name),
			attribute.String("deployment.environment", environment),
		),
	)
}

// createOTLPExporter создает gRPC экспортер для OTLP коллектора
func createOTLPExporter(ctx context.Context, endpoint string) (*otlploggrpc.Exporter, error) {
	return otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(endpoint),
		otlploggrpc.WithInsecure(), // для разработки, в продакшене следует использовать TLS
	)
}

// parseLevel преобразует строковый уровень логирования в zapcore.Level.
func parseLevel(levelStr string) zapcore.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// fieldsFromContext извлекает из контекста стандартные поля (trace_id, request_id, user_id).
func fieldsFromContext(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)

	if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
		fields = append(fields, zap.String(string(traceIDKey), traceID))
	}

	if requestID, ok := ctx.Value(requestIDKey).(string); ok && requestID != "" {
		fields = append(fields, zap.String(string(requestIDKey), requestID))
	}

	if userID, ok := ctx.Value(userIDKey).(string); ok && userID != "" {
		fields = append(fields, zap.String(string(userIDKey), userID))
	}

	return fields
}
