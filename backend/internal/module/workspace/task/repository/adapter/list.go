package adapter

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// List возвращает список задач с фильтром и пагинацией. total — общее количество без LIMIT.
func (r *Repository) List(ctx context.Context, tx *sqlx.Tx, filter *model.TaskListFilter, pagination *model.TaskPagination) ([]*model.Task, int, error) {
	return r.taskReader.List(ctx, tx, filter, pagination)
}
