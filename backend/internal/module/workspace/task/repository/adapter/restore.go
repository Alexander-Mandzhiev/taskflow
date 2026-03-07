package adapter

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Restore снимает пометку удаления. При отсутствии — model.ErrTaskNotFound.
func (r *Repository) Restore(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) error {
	return r.taskWriter.Restore(ctx, tx, taskID)
}
