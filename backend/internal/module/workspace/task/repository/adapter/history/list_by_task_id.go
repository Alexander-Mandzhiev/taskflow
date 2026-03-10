package history

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// ListHistoryByTaskID возвращает историю изменений задачи по task_id.
func (r *Adapter) ListHistoryByTaskID(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) ([]model.TaskHistory, error) {
	return r.historyReader.ListByTaskID(ctx, tx, taskID)
}
