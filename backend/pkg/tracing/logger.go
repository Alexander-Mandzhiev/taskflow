package tracing

import (
	"context"

	"go.uber.org/zap"
)

// Logger интерфейс для логирования в пакете трассировки
type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}
