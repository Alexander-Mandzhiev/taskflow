package txmanager

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

// WithTx выполняет fn в транзакции. При успехе — commit, при ошибке — rollback.
// Hooks регистрируются адаптерами через HookRegistry в context и выполняются только после успешного коммита.
// Без opts используется уровень изоляции по умолчанию драйвера (для MySQL обычно Repeatable Read).
//
// Примеры:
//
//	WithTx(ctx, fn)
//	WithTx(ctx, fn, &sql.TxOptions{Isolation: sql.LevelSerializable})
func (m *Manager) WithTx(ctx context.Context, fn func(context.Context, *sqlx.Tx) error, opts ...*sql.TxOptions) error {
	var txOpts *sql.TxOptions
	if len(opts) > 0 {
		txOpts = opts[0]
	}

	result, err := m.executeTransaction(ctx, fn, txOpts)
	if err != nil {
		return err
	}

	hooks := result.registry.GetHooks()
	if len(hooks) > 0 {
		logger.Debug(ctx, "Executing post-commit hooks", zap.Int("total_hooks_count", len(hooks)))
		m.executeHooks(ctx, hooks)
		logger.Debug(ctx, "Post-commit hooks completed", zap.Int("total_hooks_count", len(hooks)))
	}

	return nil
}
