package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// TaskRepository — единая точка доступа к данным задач и истории (контракт для адаптера).
// tx: при tx != nil все операции в транзакции; при tx == nil — вне транзакции. ID передаются как uuid.UUID.
type TaskRepository interface {
	// Create создаёт задачу. teamID и createdBy — в сигнатуре.
	Create(ctx context.Context, tx *sqlx.Tx, teamID uuid.UUID, input *model.TaskInput, createdBy uuid.UUID) (*model.Task, error)

	// GetByID возвращает задачу по id (без удалённых). При отсутствии — (nil, model.ErrTaskNotFound).
	GetByID(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) (*model.Task, error)

	// Update полностью обновляет изменяемые поля. taskID в сигнатуре. Запись в task_history — зона ответственности сервиса.
	Update(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID, input *model.TaskInput) error

	// SoftDelete помечает задачу удалённой (deleted_at = now()). При отсутствии задачи — model.ErrTaskNotFound.
	SoftDelete(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) error

	// Restore снимает пометку удаления (deleted_at = null). При отсутствии задачи — model.ErrTaskNotFound.
	Restore(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) error

	// List возвращает список задач с фильтром и пагинацией. total — общее количество без LIMIT.
	List(ctx context.Context, tx *sqlx.Tx, filter *model.TaskListFilter, pagination *model.TaskPagination) ([]*model.Task, int, error)

	// CreateHistoryEntry добавляет запись в task_history (аудит: field_name, old_value, new_value).
	CreateHistoryEntry(ctx context.Context, tx *sqlx.Tx, entry *model.TaskHistory) error

	// ListHistoryByTaskID возвращает историю изменений задачи по task_id.
	ListHistoryByTaskID(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) ([]*model.TaskHistory, error)
}
