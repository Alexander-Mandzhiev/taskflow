package task

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// Update обновляет задачу. Запись в task_history — зона ответственности сервиса.
func (r *Adapter) Update(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID, input *model.TaskInput) error {
	return r.taskWriter.Update(ctx, tx, taskID, input)
}
