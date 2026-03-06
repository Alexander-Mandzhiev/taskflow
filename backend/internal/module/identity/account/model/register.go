package model

// RegisterInput — входные данные для регистрации (из API).
type RegisterInput struct {
	Email    string
	Password string //nolint:gosec // G117: вход от клиента, не храним; только хеш в БД
	Name     string
}
