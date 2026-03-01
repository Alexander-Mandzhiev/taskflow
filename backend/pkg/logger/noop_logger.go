package logger

import (
	"context"

	"go.uber.org/zap"
)

type NoopLogger struct{}

func (l *NoopLogger) Debug(ctx context.Context, msg string, fields ...zap.Field) {}
func (l *NoopLogger) Info(ctx context.Context, msg string, fields ...zap.Field)  {}
func (l *NoopLogger) Warn(ctx context.Context, msg string, fields ...zap.Field)  {}
func (l *NoopLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {}
