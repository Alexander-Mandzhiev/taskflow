package comment

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// ListCommentsByTaskID возвращает комментарии к задаче по task_id.
func (r *Adapter) ListCommentsByTaskID(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) ([]*model.TaskComment, error) {
	return r.commentReader.ListByTaskID(ctx, tx, taskID)
}
