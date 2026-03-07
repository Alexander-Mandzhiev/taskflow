package dto

// InviteResponse — ответ на создание приглашения (запись в team_invitations). Отправка письма — отдельный сервис.
type InviteResponse struct {
	Success    bool               `json:"success"`
	Message    string             `json:"message"`
	Invitation InvitationResponse `json:"invitation,omitempty"`
}

// InvitationResponse — данные приглашения для ответа API.
type InvitationResponse struct {
	ID        string `json:"id"`
	TeamID    string `json:"team_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	ExpiresAt string `json:"expires_at"`
}
