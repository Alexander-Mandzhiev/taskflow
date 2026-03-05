package txmanager

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

// WithSerializableTx выполняет fn в транзакции с уровнем изоляции Serializable.
// При serialization failure автоматически повторяет транзакцию (до defaultMaxRetries попыток)
// с экспоненциальным backoff.
func (m *Manager) WithSerializableTx(ctx context.Context, fn func(context.Context, *sqlx.Tx) error) error {
	opts := &sql.TxOptions{Isolation: sql.LevelSerializable}

	var lastErr error
	for attempt := range defaultMaxRetries + 1 {
		if attempt > 0 {
			backoff := time.Duration(attempt*attempt) * 50 * time.Millisecond
			logger.Warn(ctx, "Retrying serializable transaction",
				zap.Int("attempt", attempt),
				zap.Duration("backoff", backoff),
				zap.Error(lastErr),
			)

			select {
			case <-ctx.Done():
				return fmt.Errorf("context cancelled during serializable retry: %w", ctx.Err())
			case <-time.After(backoff):
			}
		}

		lastErr = m.WithTx(ctx, fn, opts)
		if lastErr == nil {
			if attempt > 0 {
				if span := trace.SpanFromContext(ctx); span.IsRecording() {
					span.SetAttributes(attribute.Int("db.serializable.retries", attempt))
				}
			}
			return nil
		}

		if !isSerializationError(lastErr) {
			return lastErr
		}
	}

	return fmt.Errorf("serializable transaction failed after %d retries: %w", defaultMaxRetries, lastErr)
}
