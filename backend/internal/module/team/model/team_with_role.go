package model

// TeamWithRole — команда с ролью текущего пользователя в ней.
// Используется для списка «мои команды» (GET /api/v1/teams): одна запись = одна команда + роль пользователя.
type TeamWithRole struct {
	Team
	Role string
}
