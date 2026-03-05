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
	GetByID(ctx context.Context, tx *sqlx.Tx, id string) (*model.User, error)
	GetByEmail(ctx context.Context, tx *sqlx.Tx, email string) (*model.User, error)
}

// UserWriterRepository предоставляет методы для записи пользователей в БД.
// Мутации всегда внутри транзакции: вызывающий передаёт tx из txmanager.WithTx.
type UserWriterRepository interface {
	Create(ctx context.Context, tx *sqlx.Tx, input *model.UserInput, passwordHash string) (*model.User, error)
	Update(ctx context.Context, tx *sqlx.Tx, id string, input *model.UserInput) (*model.User, error)
	UpdatePasswordHash(ctx context.Context, tx *sqlx.Tx, id, passwordHash string) error
	Delete(ctx context.Context, tx *sqlx.Tx, id string) error
}

// UserCacheRepository предоставляет методы для кеша пользователей (по id).
// Get — чтение из кеша, Set — запись после чтения из БД, Delete — инвалидация при записи в БД.
type UserCacheRepository interface {
	Get(ctx context.Context, id string) (*model.User, error)
	Set(ctx context.Context, id string, user *model.User) error
	Delete(ctx context.Context, id string) error
}

// =============================================================================
// АДАПТЕР (объединяет reader + writer + cache)
// =============================================================================

// UserRepository — единая точка доступа к данным пользователя.
// tx — из txmanager.WithTx (*sqlx.Tx): при tx != nil чтение/запись идут в БД в транзакции; при tx == nil чтение — кеш, при промахе БД и запись в кеш, запись — БД + post-commit hook инвалидации кеша.
type UserRepository interface {
	GetByID(ctx context.Context, tx *sqlx.Tx, id string) (*model.User, error)
	GetByEmail(ctx context.Context, tx *sqlx.Tx, email string) (*model.User, error)
	Create(ctx context.Context, tx *sqlx.Tx, input *model.UserInput, passwordHash string) (*model.User, error)
	Update(ctx context.Context, tx *sqlx.Tx, id string, input *model.UserInput) (*model.User, error)
	UpdatePasswordHash(ctx context.Context, tx *sqlx.Tx, id, passwordHash string) error
	Delete(ctx context.Context, tx *sqlx.Tx, id string) error
}
