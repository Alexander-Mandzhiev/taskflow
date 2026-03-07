package converter

import (
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/team/v1/dto"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

// CreateTeamRequestToDomain конвертирует DTO запроса создания команды в доменную модель.
func CreateTeamRequestToDomain(req dto.CreateTeamRequest) *model.TeamInput {
	return &model.TeamInput{Name: req.Name}
}
