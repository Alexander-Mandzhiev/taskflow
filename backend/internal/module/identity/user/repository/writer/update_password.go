package writer

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
)

// UpdatePasswordHash обновляет хеш пароля пользователя по ID.
// Вызывается только внутри txmanager.WithTx, tx передаётся явно.
func (r *repository) UpdatePasswordHash(ctx context.Context, tx *sqlx.Tx, id, passwordHash string) error {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Question).
		Update("users").
		Set("password_hash", passwordHash).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id}).
		Where(sq.Expr("deleted_at IS NULL")).
		ToSql()
	if err != nil {
		return fmt.Errorf("build update password query: %w", err)
	}

	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update password exec: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return model.ErrUserNotFound
	}
	return nil
}
