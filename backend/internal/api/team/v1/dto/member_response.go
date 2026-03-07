package dto

// MemberResponse — участник команды в ответе API.
type MemberResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	TeamID    string `json:"team_id"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}
