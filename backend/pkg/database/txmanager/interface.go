package txmanager

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// PostCommitHook выполняется после успешного коммита транзакции.
// Ошибки в hooks логируются, но не прерывают выполнение.
type PostCommitHook func(context.Context) error

// TxManager — интерфейс для управления транзакциями (удобно для моков в тестах).
// Колбек получает *sqlx.Tx — один тип по всему стеку (репо, сервисы), без обёрток sql.Tx → sqlx.Tx в репозитории.
type TxManager interface {
	// WithTx выполняет функцию в транзакции с опциональной поддержкой post-commit hooks.
	// Hooks регистрируются адаптерами через HookRegistry в context и выполняются только после успешного коммита.
	// Без opts используется уровень изоляции по умолчанию драйвера (для MySQL обычно Repeatable Read).
	WithTx(ctx context.Context, fn func(context.Context, *sqlx.Tx) error, opts ...*sql.TxOptions) error
	// WithSerializableTx выполняет функцию в транзакции с уровнем изоляции Serializable.
	WithSerializableTx(ctx context.Context, fn func(context.Context, *sqlx.Tx) error) error
}
