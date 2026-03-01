package dto

// RegisterResponse — ответ на регистрацию.
type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
