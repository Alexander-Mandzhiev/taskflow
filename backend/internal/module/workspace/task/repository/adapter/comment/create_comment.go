package comment

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// CreateComment создаёт комментарий к задаче (task_id, user_id, content).
func (r *Adapter) CreateComment(ctx context.Context, tx *sqlx.Tx, taskID, userID uuid.UUID, content string) (*model.TaskComment, error) {
	return r.commentWriter.Create(ctx, tx, taskID, userID, content)
}
