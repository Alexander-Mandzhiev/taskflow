package task

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Restore снимает пометку удаления. При отсутствии — model.ErrTaskNotFound.
func (r *Adapter) Restore(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) error {
	return r.taskWriter.Restore(ctx, tx, taskID)
}
