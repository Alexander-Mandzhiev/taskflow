package dto

// LoginResponse — ответ на вход. Токены отдаются только в cookie (access_token, refresh_token).
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
