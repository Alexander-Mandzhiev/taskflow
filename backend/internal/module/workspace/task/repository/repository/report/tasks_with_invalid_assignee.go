package report

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/converter"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/resources"
)

// TasksWithInvalidAssignee возвращает задачи, у которых assignee не является участником команды задачи.
func (r *repository) TasksWithInvalidAssignee(ctx context.Context, tx *sqlx.Tx) ([]*model.Task, error) {
	const query = `
		SELECT t.id, t.title, t.description, t.status, t.assignee_id, t.team_id, t.created_by, t.created_at, t.updated_at, t.completed_at, t.deleted_at
		FROM tasks t
		WHERE t.deleted_at IS NULL AND t.assignee_id IS NOT NULL
		AND NOT EXISTS (
			SELECT 1 FROM team_members tm WHERE tm.team_id = t.team_id AND tm.user_id = t.assignee_id
		)
		ORDER BY t.updated_at DESC
	`
	var rows []resources.TaskRow
	if tx != nil {
		if err := tx.SelectContext(ctx, &rows, query); err != nil {
			return nil, toDomainError(err)
		}
	} else {
		if err := r.readPool.SelectContext(ctx, &rows, query); err != nil {
			return nil, toDomainError(err)
		}
	}

	out := make([]*model.Task, 0, len(rows))
	for i := range rows {
		task, err := converter.ToDomainTask(rows[i])
		if err != nil {
			return nil, fmt.Errorf("convert row %d: %w", i, err)
		}
		out = append(out, &task)
	}
	return out, nil
}
