package reader

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
)

// GetByID возвращает пользователя по ID (без удалённых).
// При tx != nil запрос выполняется в транзакции.
func (r *repository) GetByID(ctx context.Context, tx *sqlx.Tx, id string) (*model.User, error) {
	return r.getOne(ctx, tx, sq.Eq{"id": id}, "get by id")
}
