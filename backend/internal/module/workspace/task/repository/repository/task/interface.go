package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// TaskReaderRepository — чтение из таблицы tasks.
// tx != nil — запрос в транзакции; tx == nil — вне транзакции.
type TaskReaderRepository interface {
	// GetByID возвращает задачу по id (без удалённых). При отсутствии — (nil, model.ErrTaskNotFound).
	GetByID(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) (*model.Task, error)

	// GetByIDIncludeDeleted возвращает задачу по id в том числе удалённую. При отсутствии — (nil, model.ErrTaskNotFound).
	GetByIDIncludeDeleted(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) (*model.Task, error)

	// List возвращает список задач с фильтром и пагинацией. total — общее количество записей без учёта LIMIT.
	List(ctx context.Context, tx *sqlx.Tx, filter *model.TaskListFilter, pagination *model.TaskPagination) ([]*model.Task, int, error)
}

// TaskWriterRepository — запись в таблицу tasks. Мутации в транзакции (tx из txmanager.WithTx).
type TaskWriterRepository interface {
	// Create создаёт задачу. teamID и createdBy — в сигнатуре.
	Create(ctx context.Context, tx *sqlx.Tx, teamID uuid.UUID, input *model.TaskInput, createdBy uuid.UUID) (*model.Task, error)

	// Update полностью обновляет изменяемые поля (title, description, status, assignee_id). taskID в сигнатуре.
	Update(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID, input *model.TaskInput) error

	// SoftDelete помечает задачу удалённой (deleted_at = now()). При отсутствии — model.ErrTaskNotFound.
	SoftDelete(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) error

	// Restore снимает пометку удаления (deleted_at = null). При отсутствии — model.ErrTaskNotFound.
	Restore(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) error
}
