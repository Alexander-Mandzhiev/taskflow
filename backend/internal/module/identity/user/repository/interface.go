package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
)

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
