package logger

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/sdk/log"
	"go.uber.org/zap"
)

// Глобальные переменные пакета
var (
	globalLogger *logger             // глобальный экземпляр логгера
	initOnce     sync.Once           // обеспечивает единократную инициализацию
	level        zap.AtomicLevel     // уровень логирования (может изменяться динамически)
	otelProvider *log.LoggerProvider // OTLP provider для graceful shutdown
)

// SetLevel изменяет уровень логирования у уже инициализированного глобального логгера.
func SetLevel(levelStr string) {
	if level == (zap.AtomicLevel{}) {
		return
	}

	level.SetLevel(parseLevel(levelStr))
}

// Logger возвращает глобальный логгер.
func Logger() *logger {
	return globalLogger
}

// Sync сбрасывает буферы логгера.
func Sync() error {
	if globalLogger != nil {
		return globalLogger.zapLogger.Sync()
	}

	return nil
}

// With создает новый enrich-aware логгер с дополнительными полями
func With(fields ...zap.Field) *logger {
	if globalLogger == nil {
		return &logger{zapLogger: zap.NewNop()}
	}

	return &logger{
		zapLogger: globalLogger.zapLogger.With(fields...),
	}
}

// WithContext добавляет к логгеру поля из контекста (trace_id/request_id/user_id).
func WithContext(ctx context.Context) *logger {
	if globalLogger == nil {
		return &logger{zapLogger: zap.NewNop()}
	}

	return &logger{
		zapLogger: globalLogger.zapLogger.With(fieldsFromContext(ctx)...),
	}
}

// Debug пишет сообщение уровня DEBUG с полями.
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Debug(ctx, msg, fields...)
	}
}

// Info пишет сообщение уровня INFO с полями.
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Info(ctx, msg, fields...)
	}
}

// Warn пишет сообщение уровня WARN с полями.
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Warn(ctx, msg, fields...)
	}
}

// Error пишет сообщение уровня ERROR с полями.
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Error(ctx, msg, fields...)
	}
}

// Fatal пишет сообщение уровня FATAL и завершает процесс.
func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Fatal(ctx, msg, fields...)
	}
}

// WithIDs добавляет request/trace id в контекст, не генерируя их.
func WithIDs(ctx context.Context, traceID, requestID string) context.Context {
	if traceID != "" {
		ctx = context.WithValue(ctx, traceIDKey, traceID)
	}
	if requestID != "" {
		ctx = context.WithValue(ctx, requestIDKey, requestID)
	}
	return ctx
}

// Получить trace_id из контекста
func TraceIDFrom(ctx context.Context) string {
	if v, ok := ctx.Value(traceIDKey).(string); ok {
		return v
	}
	return ""
}

// Получить request_id из контекста
func RequestIDFrom(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey).(string); ok {
		return v
	}
	return ""
}
