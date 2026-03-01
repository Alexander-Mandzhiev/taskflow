package dto

import "github.com/google/uuid"

// RegisterResponse — ответ на регистрацию.
type RegisterResponse struct {
	UserID  uuid.UUID `json:"user_id"`
	Message string    `json:"message"`
}
