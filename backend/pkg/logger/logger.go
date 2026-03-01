package logger

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/sdk/log"
	"go.uber.org/zap"

	"mkk/pkg/ctxkey"
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
	if globalLogger == nil {
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
		allFields := append(fieldsFromContext(ctx), fields...)
		globalLogger.zapLogger.Debug(msg, allFields...)
	}
}

// Info пишет сообщение уровня INFO с полями.
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	if globalLogger != nil {
		allFields := append(fieldsFromContext(ctx), fields...)
		globalLogger.zapLogger.Info(msg, allFields...)
	}
}

// Warn пишет сообщение уровня WARN с полями.
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	if globalLogger != nil {
		allFields := append(fieldsFromContext(ctx), fields...)
		globalLogger.zapLogger.Warn(msg, allFields...)
	}
}

// Error пишет сообщение уровня ERROR с полями.
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	if globalLogger != nil {
		allFields := append(fieldsFromContext(ctx), fields...)
		globalLogger.zapLogger.Error(msg, allFields...)
	}
}

// Fatal пишет сообщение уровня FATAL и завершает процесс.
func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	if globalLogger != nil {
		allFields := append(fieldsFromContext(ctx), fields...)
		globalLogger.zapLogger.Fatal(msg, allFields...)
	}
}

// WithIDs добавляет request/trace id в контекст, не генерируя их.
func WithIDs(ctx context.Context, traceID, requestID string) context.Context {
	if traceID != "" {
		ctx = context.WithValue(ctx, ctxkey.TraceID, traceID)
	}
	if requestID != "" {
		ctx = context.WithValue(ctx, ctxkey.RequestID, requestID)
	}
	return ctx
}

// Получить trace_id из контекста
func TraceIDFrom(ctx context.Context) string {
	if v, ok := ctx.Value(ctxkey.TraceID).(string); ok {
		return v
	}
	return ""
}

// Получить request_id из контекста
func RequestIDFrom(ctx context.Context) string {
	if v, ok := ctx.Value(ctxkey.RequestID).(string); ok {
		return v
	}
	return ""
}
