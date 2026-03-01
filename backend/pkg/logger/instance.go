package logger

import (
	"context"

	"go.uber.org/zap"
)

// logger представляет экземпляр логгера
type logger struct {
	zapLogger *zap.Logger
}

// Instance methods для enrich loggers (logger)

func (l *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Debug(msg, allFields...)
}

// Info пишет сообщение INFO, добавляя поля из контекста.
func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Info(msg, allFields...)
}

// Warn пишет сообщение WARN, добавляя поля из контекста.
func (l *logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Warn(msg, allFields...)
}

// Error пишет сообщение ERROR, добавляя поля из контекста.
func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Error(msg, allFields...)
}

// Fatal пишет сообщение FATAL, добавляя поля из контекста, и завершает процесс.
func (l *logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Fatal(msg, allFields...)
}
