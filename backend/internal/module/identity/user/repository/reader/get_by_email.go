package reader

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"mkk/internal/module/identity/user/model"
)

// GetByEmail возвращает пользователя по email (без удалённых). При tx != nil запрос в транзакции.
func (r *repository) GetByEmail(ctx context.Context, tx *sqlx.Tx, email string) (*model.User, error) {
	return r.getOne(ctx, tx, sq.Eq{"email": email}, "get by email")
}
