package model

// TeamWithMembers — команда с составом участников (обогащённая модель для ответа фронту).
// Используется для детальной страницы команды (GET /api/v1/teams/{id}): один вызов репозитория возвращает команду + членов.
type TeamWithMembers struct {
	Team
	Members []*TeamMember
}
