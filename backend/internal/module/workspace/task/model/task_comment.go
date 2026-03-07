package model

import (
	"time"

	"github.com/google/uuid"
)

// TaskComment — комментарий к задаче (таблица task_comments).
// API комментариев по ТЗ не перечислен — модель для целостности БД и на будущее.
type TaskComment struct {
	ID        uuid.UUID
	TaskID    uuid.UUID
	UserID    uuid.UUID
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
