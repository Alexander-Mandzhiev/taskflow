package writer

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
)

// Delete выполняет мягкое удаление пользователя по ID.
// Вызывается только внутри txmanager.WithTx, tx передаётся явно.
// Валидация id (формат UUID) выполняется в сервисном и API-слое.
func (r *repository) Delete(ctx context.Context, tx *sqlx.Tx, id string) error {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Question).
		Update("users").
		Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id}).
		Where(sq.Expr("deleted_at IS NULL")).
		ToSql()
	if err != nil {
		return fmt.Errorf("build delete query: %w", err)
	}
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return toDomainError(err)
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
