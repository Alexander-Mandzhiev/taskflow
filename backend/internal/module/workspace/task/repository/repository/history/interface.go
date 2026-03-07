package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// TaskHistoryReaderRepository — чтение из таблицы task_history.
type TaskHistoryReaderRepository interface {
	// ListByTaskID возвращает историю изменений задачи по task_id, упорядоченную по changed_at.
	ListByTaskID(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) ([]*model.TaskHistory, error)
}

// TaskHistoryWriterRepository — запись в таблицу task_history.
type TaskHistoryWriterRepository interface {
	// Create создаёт запись истории (id, changed_at могут генерироваться в БД или задаваться вызывающим).
	Create(ctx context.Context, tx *sqlx.Tx, entry *model.TaskHistory) error
}
