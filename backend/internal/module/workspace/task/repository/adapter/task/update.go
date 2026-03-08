package task

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// Update обновляет задачу. Запись в task_history — зона ответственности сервиса.
// После успешного обновления регистрирует post-commit хук инвалидации кеша списка по команде.
func (r *Adapter) Update(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID, input *model.TaskInput) error {
	current, err := r.taskReader.GetByID(ctx, tx, taskID)
	if err != nil {
		return err
	}
	if err := r.taskWriter.Update(ctx, tx, taskID, input); err != nil {
		return err
	}
	r.registerInvalidateHook(ctx, current.TeamID)
	return nil
}
