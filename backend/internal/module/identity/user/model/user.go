package model

import (
	"time"

	"github.com/google/uuid"
)

// User — модель пользователя для чтения (все поля, включая служебные).
type User struct {
	ID           uuid.UUID
	Email        string
	Name         string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
