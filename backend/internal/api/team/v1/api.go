package team_v1

import (
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/service"
)

// API — HTTP-слой команд: создание, список, получение по id, приглашение.
// Требует JWT (user_id из контекста).
type API struct {
	teamService service.TeamService
}

// NewAPI создаёт API.
func NewAPI(teamService service.TeamService) *API {
	return &API{teamService: teamService}
}
