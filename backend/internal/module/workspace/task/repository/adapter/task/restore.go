package task

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Restore снимает пометку удаления. При отсутствии — model.ErrTaskNotFound.
// После успешного восстановления регистрирует post-commit хук инвалидации кеша списка по команде.
func (r *Adapter) Restore(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) error {
	task, err := r.taskReader.GetByIDIncludeDeleted(ctx, tx, taskID)
	if err != nil {
		return err
	}
	if err := r.taskWriter.Restore(ctx, tx, taskID); err != nil {
		return err
	}
	r.registerInvalidateHook(ctx, task.TeamID)
	return nil
}
