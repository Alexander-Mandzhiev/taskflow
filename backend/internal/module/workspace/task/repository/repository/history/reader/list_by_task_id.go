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

// ListByTaskID возвращает историю изменений задачи по task_id, упорядоченную по changed_at.
func (r *repository) ListByTaskID(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) ([]*model.TaskHistory, error) {
	const query = `
		SELECT id, task_id, changed_by, field_name, old_value, new_value, changed_at
		FROM task_history WHERE task_id = ? ORDER BY changed_at ASC
	`
	var rows []resources.TaskHistoryRow
	if tx != nil {
		if err := tx.SelectContext(ctx, &rows, query, taskID.String()); err != nil {
			return nil, toDomainError(err)
		}
	} else {
		if err := r.readPool.SelectContext(ctx, &rows, query, taskID.String()); err != nil {
			return nil, toDomainError(err)
		}
	}

	out := make([]*model.TaskHistory, 0, len(rows))
	for i := range rows {
		h, err := converter.ToDomainTaskHistory(rows[i])
		if err != nil {
			return nil, fmt.Errorf("convert row %d: %w", i, err)
		}
		out = append(out, &h)
	}
	return out, nil
}
