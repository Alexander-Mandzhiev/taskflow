package dto

// InviteResponse — ответ на приглашение в команду.
type InviteResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Member  MemberResponse `json:"member,omitempty"`
}
