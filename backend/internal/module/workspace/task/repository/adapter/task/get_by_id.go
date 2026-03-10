package task

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// GetByID возвращает задачу по id (без удалённых). При отсутствии — model.ErrTaskNotFound.
func (r *Adapter) GetByID(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) (model.Task, error) {
	return r.taskReader.GetByID(ctx, tx, taskID)
}
