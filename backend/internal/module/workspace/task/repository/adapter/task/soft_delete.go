package task

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// SoftDelete помечает задачу удалённой. При отсутствии — model.ErrTaskNotFound.
// После успешного удаления регистрирует post-commit хук инвалидации кеша списка по команде.
func (r *Adapter) SoftDelete(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) error {
	task, err := r.taskReader.GetByID(ctx, tx, taskID)
	if err != nil {
		return err
	}
	if err := r.taskWriter.SoftDelete(ctx, tx, taskID); err != nil {
		return err
	}
	r.registerInvalidateHook(ctx, task.TeamID)
	return nil
}
