package model

import (
	"time"

	"github.com/google/uuid"
)

// Team — модель команды (таблица teams).
type Team struct {
	ID        uuid.UUID
	Name      string
	CreatedBy uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
