package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// TaskCommentReaderRepository — чтение из таблицы task_comments.
type TaskCommentReaderRepository interface {
	// ListByTaskID возвращает комментарии к задаче по task_id, без удалённых (deleted_at IS NULL), упорядоченные по created_at.
	ListByTaskID(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) ([]model.TaskComment, error)
}

// TaskCommentWriterRepository — запись в таблицу task_comments.
type TaskCommentWriterRepository interface {
	// Create создаёт комментарий; id, created_at, updated_at генерируются в БД или задаются вызывающим.
	Create(ctx context.Context, tx *sqlx.Tx, taskID, userID uuid.UUID, content string) (model.TaskComment, error)
}
