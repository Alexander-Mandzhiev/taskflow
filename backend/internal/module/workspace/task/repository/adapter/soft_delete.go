package adapter

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// SoftDelete помечает задачу удалённой. При отсутствии — model.ErrTaskNotFound.
func (r *Repository) SoftDelete(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) error {
	return r.taskWriter.SoftDelete(ctx, tx, taskID)
}
