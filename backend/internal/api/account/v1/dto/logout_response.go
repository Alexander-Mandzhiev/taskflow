package dto

// LogoutResponse — ответ на выход.
type LogoutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
