package txmanager

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

// transactionResult результат выполнения транзакции (для вызова hooks после коммита).
type transactionResult struct {
	duration time.Duration
	registry *HookRegistry
}

func (m *Manager) executeTransaction(
	ctx context.Context,
	fn func(context.Context, *sqlx.Tx) error,
	opts *sql.TxOptions,
) (*transactionResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context already cancelled: %w", err)
	}

	var isolationLevel string
	if opts != nil {
		isolationLevel = isolationLevelToString(opts.Isolation)
	} else {
		isolationLevel = "default"
	}

	registry := NewHookRegistry()
	ctx = WithHookRegistry(ctx, registry)

	ctx, span := m.tracer.Start(ctx, "db.transaction",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("db.operation", "transaction"),
			attribute.String("db.transaction.isolation", isolationLevel),
		),
	)
	defer span.End()

	logger.Debug(ctx, "Transaction started",
		zap.String("isolation_level", isolationLevel),
	)

	start := time.Now()

	tx, err := m.db.BeginTxx(ctx, opts)
	if err != nil {
		logger.Error(ctx, "Failed to begin transaction", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(attribute.String("db.status", "error"))
		return nil, fmt.Errorf("begin transaction: %w", err)
	}

	if err := fn(ctx, tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logger.Error(ctx, "Failed to rollback transaction",
				zap.Error(rollbackErr),
				zap.NamedError("original_error", err),
			)
			span.RecordError(rollbackErr)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			span.SetAttributes(attribute.String("db.status", "rollback_failed"))
			return nil, fmt.Errorf("transaction failed: %w; rollback failed: %w", err, rollbackErr)
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(attribute.String("db.status", "rollback"))
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		logger.Error(ctx, "Failed to commit transaction", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(attribute.String("db.status", "error"))
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	duration := time.Since(start)
	span.SetStatus(codes.Ok, "success")
	span.SetAttributes(
		attribute.String("db.status", "commit"),
		attribute.Float64("db.duration_us", float64(duration.Microseconds())),
	)

	logger.Debug(ctx, "Transaction committed", zap.Duration("duration", duration))

	return &transactionResult{
		duration: duration,
		registry: registry,
	}, nil
}
