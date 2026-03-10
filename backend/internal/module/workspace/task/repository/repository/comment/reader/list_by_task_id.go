package reader

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/converter"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/resources"
)

// ListByTaskID возвращает комментарии к задаче по task_id (без удалённых), упорядоченные по created_at.
func (r *repository) ListByTaskID(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) ([]model.TaskComment, error) {
	const query = `
		SELECT id, task_id, user_id, content, created_at, updated_at, deleted_at
		FROM task_comments WHERE task_id = ? AND deleted_at IS NULL ORDER BY created_at ASC
	`
	var rows []resources.TaskCommentRow
	if tx != nil {
		if err := tx.SelectContext(ctx, &rows, query, taskID.String()); err != nil {
			return nil, toDomainError(err)
		}
	} else {
		if err := r.readPool.SelectContext(ctx, &rows, query, taskID.String()); err != nil {
			return nil, toDomainError(err)
		}
	}

	out := make([]model.TaskComment, 0, len(rows))
	for i := range rows {
		c, err := converter.ToDomainTaskComment(rows[i])
		if err != nil {
			return nil, fmt.Errorf("convert row %d: %w", i, err)
		}
		out = append(out, c)
	}
	return out, nil
}
