package writer

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/converter"
)

// Create добавляет запись в task_history. id и changed_at могут быть заданы в entry; иначе генерируются.
func (r *repository) Create(ctx context.Context, tx *sqlx.Tx, entry model.TaskHistory) error {
	row := converter.ToRepoTaskHistory(entry)
	if entry.ID == uuid.Nil {
		row.ID = uuid.New().String()
	}

	const query = `
		INSERT INTO task_history (id, task_id, changed_by, field_name, old_value, new_value, changed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	args := []interface{}{row.ID, row.TaskID, row.ChangedBy, row.FieldName, row.OldValue, row.NewValue, row.ChangedAt}

	if tx != nil {
		_, err := tx.ExecContext(ctx, query, args...)
		return toDomainError(err)
	}
	_, err := r.writePool.ExecContext(ctx, query, args...)
	return toDomainError(err)
}
