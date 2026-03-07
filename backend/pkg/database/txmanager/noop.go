package txmanager

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Noop — реализация TxManager для тестов: выполняет колбек с пустой транзакцией без реального БД.
// Подходит для юнит-тестов сервисов, где репозитории замоканы.
type Noop struct{}

// WithTx выполняет fn с пустым *sqlx.Tx.
func (Noop) WithTx(ctx context.Context, fn func(context.Context, *sqlx.Tx) error, _ ...*sql.TxOptions) error {
	return fn(ctx, &sqlx.Tx{})
}

// WithSerializableTx выполняет fn с пустым *sqlx.Tx.
func (Noop) WithSerializableTx(ctx context.Context, fn func(context.Context, *sqlx.Tx) error) error {
	return fn(ctx, &sqlx.Tx{})
}

var _ TxManager = (*Noop)(nil)
