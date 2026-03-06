package model

// LoginInput — входные данные для логина (из API + контекст запроса).
type LoginInput struct {
	Email     string
	Password  string //nolint:gosec // G117: вход от клиента, не храним; только хеш в БД
	UserAgent string
	IP        string
}
