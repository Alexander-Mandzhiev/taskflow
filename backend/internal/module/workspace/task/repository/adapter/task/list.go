package task

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// List возвращает список задач по фильтру (критерии + limit/offset в filter). total — общее количество без LIMIT.
func (r *Adapter) List(ctx context.Context, tx *sqlx.Tx, filter *model.TaskListFilter) ([]*model.Task, int, error) {
	return r.taskReader.List(ctx, tx, filter)
}
