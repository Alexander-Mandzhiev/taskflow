package report

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// TasksWithInvalidAssignee возвращает задачи с assignee не из команды.
func (r *Adapter) TasksWithInvalidAssignee(ctx context.Context, tx *sqlx.Tx) ([]model.Task, error) {
	return r.reader.TasksWithInvalidAssignee(ctx, tx)
}
