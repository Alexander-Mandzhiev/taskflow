package model

import (
	"time"

	"github.com/google/uuid"
)

// TaskHistory — одна запись истории изменений задачи (таблица task_history).
type TaskHistory struct {
	ID        uuid.UUID
	TaskID    uuid.UUID
	ChangedBy uuid.UUID
	FieldName string
	OldValue  string
	NewValue  string
	ChangedAt time.Time
}
