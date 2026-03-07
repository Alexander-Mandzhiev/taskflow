package writer

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// SoftDelete помечает задачу удалённой (deleted_at = NOW()). При отсутствии — model.ErrTaskNotFound.
func (r *repository) SoftDelete(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) error {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Question).
		Update("tasks").
		Set("deleted_at", sq.Expr("NOW()")).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": taskID.String()}).
		Where(sq.Expr("deleted_at IS NULL")).
		ToSql()
	if err != nil {
		return fmt.Errorf("build soft delete query: %w", err)
	}

	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return toDomainError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return model.ErrTaskNotFound
	}
	return nil
}
