package reader

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository/converter"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository/resources"
)

// getOne выполняет SELECT по условию where и возвращает одну запись или model.ErrUserNotFound.
// При tx != nil запрос в транзакции, иначе через readPool.
func (r *repository) getOne(ctx context.Context, tx *sqlx.Tx, where sq.Eq, errLabel string) (model.User, error) {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Question).
		Select("id", "email", "name", "password_hash", "created_at", "updated_at", "deleted_at").
		From("users").
		Where(where).
		Where(sq.Expr("deleted_at IS NULL")).
		Limit(1).
		ToSql()
	if err != nil {
		return model.User{}, fmt.Errorf("build %s query: %w", errLabel, err)
	}

	var row resources.UserRow
	if tx != nil {
		err = tx.GetContext(ctx, &row, query, args...)
	} else {
		err = r.readPool.GetContext(ctx, &row, query, args...)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, model.ErrUserNotFound
		}
		return model.User{}, toDomainError(err)
	}

	user, err := converter.ToDomainUser(row)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
