package txmanager

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
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
		m.executeHooks(ctx, hooks)
	}

	return nil
}
