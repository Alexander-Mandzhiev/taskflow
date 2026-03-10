package writer

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/converter"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/resources"
)

// selectByID читает комментарий по ID после Create (в той же транзакции или на том же пуле).
func (r *repository) selectByID(ctx context.Context, tx *sqlx.Tx, commentID string) (model.TaskComment, error) {
	const query = `
		SELECT id, task_id, user_id, content, created_at, updated_at, deleted_at
		FROM task_comments WHERE id = ? AND deleted_at IS NULL LIMIT 1
	`
	var row resources.TaskCommentRow
	var err error
	if tx != nil {
		err = tx.GetContext(ctx, &row, query, commentID)
	} else {
		err = r.writePool.GetContext(ctx, &row, query, commentID)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.TaskComment{}, model.ErrCommentNotFound
		}
		return model.TaskComment{}, toDomainError(err)
	}
	comment, err := converter.ToDomainTaskComment(row)
	if err != nil {
		return model.TaskComment{}, err
	}
	return comment, nil
}
