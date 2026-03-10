package task

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// GetByIDIncludeDeleted возвращает задачу по id в том числе удалённую (для Restore).
func (r *Adapter) GetByIDIncludeDeleted(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) (model.Task, error) {
	return r.taskReader.GetByIDIncludeDeleted(ctx, tx, taskID)
}
