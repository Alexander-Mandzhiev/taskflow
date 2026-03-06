package writer

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

// selectByID читает полную строку пользователя по ID внутри текущей транзакции.
// Используется в Create и Update для возврата сохранённой сущности (MySQL не поддерживает RETURNING).
func (r *repository) selectByID(ctx context.Context, tx *sqlx.Tx, id string) (*model.User, error) {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Question).
		Select("id", "email", "name", "password_hash", "created_at", "updated_at", "deleted_at").
		From("users").
		Where(sq.Eq{"id": id}).
		Where(sq.Expr("deleted_at IS NULL")).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select query: %w", err)
	}

	var row resources.UserRow
	if err := tx.GetContext(ctx, &row, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("select after write: %w", err)
	}

	user, err := converter.ToDomainUser(row)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
