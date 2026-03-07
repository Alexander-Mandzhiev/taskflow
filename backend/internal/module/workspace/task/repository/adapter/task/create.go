package task

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// Create создаёт задачу. teamID и createdBy — в сигнатуре.
func (r *Adapter) Create(ctx context.Context, tx *sqlx.Tx, teamID uuid.UUID, input *model.TaskInput, createdBy uuid.UUID) (*model.Task, error) {
	return r.taskWriter.Create(ctx, tx, teamID, input, createdBy)
}
