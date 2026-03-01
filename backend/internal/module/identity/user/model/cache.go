package model

import (
	"time"
)

// UserCache — модель пользователя для хранения в кеше (GetByID).
// PasswordHash намеренно исключён: кеш хранит только публичные данные профиля.
// Для операций, требующих PasswordHash (логин, смена пароля), используется БД напрямую.
type UserCache struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
