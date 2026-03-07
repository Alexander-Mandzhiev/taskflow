package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// TaskRepository — доступ к данным задач (таблица tasks). Контракт для адаптера.
// tx: при tx != nil — в транзакции; при tx == nil — вне. ID — uuid.UUID.
type TaskRepository interface {
	Create(ctx context.Context, tx *sqlx.Tx, teamID uuid.UUID, input *model.TaskInput, createdBy uuid.UUID) (*model.Task, error)
	GetByID(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) (*model.Task, error)
	GetByIDIncludeDeleted(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) (*model.Task, error)
	Update(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID, input *model.TaskInput) error
	SoftDelete(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) error
	Restore(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) error
	List(ctx context.Context, tx *sqlx.Tx, filter *model.TaskListFilter) ([]*model.Task, int, error)
}

// TaskHistoryRepository — доступ к истории изменений задач (таблица task_history). Контракт для адаптера.
type TaskHistoryRepository interface {
	CreateHistoryEntry(ctx context.Context, tx *sqlx.Tx, entry *model.TaskHistory) error
	ListHistoryByTaskID(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) ([]*model.TaskHistory, error)
}

// TaskCommentRepository — доступ к комментариям задач (таблица task_comments). Контракт для адаптера.
type TaskCommentRepository interface {
	ListCommentsByTaskID(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) ([]*model.TaskComment, error)
	CreateComment(ctx context.Context, tx *sqlx.Tx, taskID, userID uuid.UUID, content string) (*model.TaskComment, error)
}

// ReportRepository — отчёты по задачам и командам (сложные запросы с JOIN/агрегацией). Контракт для адаптера.
type ReportRepository interface {
	TeamTaskStats(ctx context.Context, tx *sqlx.Tx, since time.Time) ([]*model.TeamTaskStats, error)
	TopCreatorsByTeam(ctx context.Context, tx *sqlx.Tx, since time.Time, limit int) ([]*model.TeamTopCreator, error)
	TasksWithInvalidAssignee(ctx context.Context, tx *sqlx.Tx) ([]*model.Task, error)
}
