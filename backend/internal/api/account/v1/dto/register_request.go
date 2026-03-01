package dto

// RegisterRequest — запрос на регистрацию.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"` //nolint:gosec // G117: поле тела запроса, не хранится
	Name     string `json:"name" validate:"omitempty,max=255"`
}
