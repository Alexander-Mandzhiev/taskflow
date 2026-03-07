package history

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// CreateHistoryEntry добавляет запись в task_history (аудит: field_name, old_value, new_value).
func (r *Adapter) CreateHistoryEntry(ctx context.Context, tx *sqlx.Tx, entry *model.TaskHistory) error {
	return r.historyWriter.Create(ctx, tx, entry)
}
