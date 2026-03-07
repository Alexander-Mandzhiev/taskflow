package dto

// InviteRequest — запрос на приглашение пользователя в команду по email.
type InviteRequest struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=admin member"`
}
