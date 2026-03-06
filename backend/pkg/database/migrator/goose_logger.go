package migrator

import (
	"context"

	"go.uber.org/zap"
)

// Logger — интерфейс для логирования в пакете миграций (goose).
type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}
