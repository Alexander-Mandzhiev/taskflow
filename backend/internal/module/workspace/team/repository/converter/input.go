package converter

import (
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/resources"
)

// ToRepoTeamInput преобразует доменный TeamInput в ресурс репозитория.
// Вызывающий гарантирует m != nil; проверка на nil выполняется в сервисе до вызова репозитория.
func ToRepoTeamInput(m *model.TeamInput) resources.TeamInput {
	if m == nil {
		return resources.TeamInput{}
	}
	return resources.TeamInput{
		Name: m.Name,
	}
}
