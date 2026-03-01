package dto

// LoginResponse — ответ на вход.
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
