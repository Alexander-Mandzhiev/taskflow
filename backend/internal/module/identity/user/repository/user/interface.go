package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
)

// =============================================================================
// ПОЛЬЗОВАТЕЛЬ (User)
// =============================================================================

// UserReaderRepository предоставляет методы для чтения пользователей из БД.
// tx — опционально: при tx != nil запрос выполняется в транзакции (для валидации внутри мутаций).
type UserReaderRepository interface {
	GetByID(ctx context.Context, tx *sqlx.Tx, id string) (model.User, error)
	GetByEmail(ctx context.Context, tx *sqlx.Tx, email string) (model.User, error)
}

// UserWriterRepository предоставляет методы для записи пользователей в БД.
// Мутации всегда внутри транзакции: вызывающий передаёт tx из txmanager.WithTx.
type UserWriterRepository interface {
	Create(ctx context.Context, tx *sqlx.Tx, input model.UserInput, passwordHash string) (model.User, error)
	Update(ctx context.Context, tx *sqlx.Tx, id string, input model.UserInput) (model.User, error)
	UpdatePasswordHash(ctx context.Context, tx *sqlx.Tx, id, passwordHash string) error
	Delete(ctx context.Context, tx *sqlx.Tx, id string) error
}
