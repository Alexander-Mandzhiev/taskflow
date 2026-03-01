package adapter

import (
	"context"

	"github.com/jmoiron/sqlx"

	"mkk/internal/module/identity/user/model"
)

// GetByEmail возвращает пользователя по email. Всегда читает из БД (кеш по email не ведётся).
// Используется для аутентификации: возвращает полную модель с PasswordHash.
func (r *Repository) GetByEmail(ctx context.Context, tx *sqlx.Tx, email string) (*model.User, error) {
	return r.reader.GetByEmail(ctx, tx, email)
}
