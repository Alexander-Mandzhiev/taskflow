package model

import (
	"time"

	"github.com/google/uuid"
)

// Статусы задачи (tasks.status).
const (
	TaskStatusTodo       = "todo"
	TaskStatusInProgress = "in_progress"
	TaskStatusDone       = "done"
)

// IsValidTaskStatus возвращает true, если статус один из допустимых (todo, in_progress, done).
func IsValidTaskStatus(s string) bool {
	switch s {
	case TaskStatusTodo, TaskStatusInProgress, TaskStatusDone:
		return true
	default:
		return false
	}
}

// Task — модель задачи (таблица tasks).
type Task struct {
	ID          uuid.UUID
	Title       string
	Description string
	Status      string
	AssigneeID  *uuid.UUID
	TeamID      uuid.UUID
	CreatedBy   uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt *time.Time // когда задача переведена в статус done; NULL если не done или снята с done
	DeletedAt   *time.Time
}
