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

// selectByID читает задачу по ID после Create (в той же транзакции или на том же пуле).
func (r *repository) selectByID(ctx context.Context, tx *sqlx.Tx, taskID string) (*model.Task, error) {
	const query = `
		SELECT id, title, description, status, assignee_id, team_id, created_by, created_at, updated_at, completed_at, deleted_at
		FROM tasks WHERE id = ? AND deleted_at IS NULL LIMIT 1
	`
	var row resources.TaskRow
	var err error
	if tx != nil {
		err = tx.GetContext(ctx, &row, query, taskID)
	} else {
		err = r.writePool.GetContext(ctx, &row, query, taskID)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrTaskNotFound
		}
		return nil, toDomainError(err)
	}
	task, err := converter.ToDomainTask(row)
	if err != nil {
		return nil, err
	}
	return &task, nil
}
