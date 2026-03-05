package dto

// LoginRequest — запрос на вход.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"` //nolint:gosec // G117: поле тела запроса, не хранится
}
