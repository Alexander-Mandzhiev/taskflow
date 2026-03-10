package converter

import (
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository/resources"
)

// ToRepoInput преобразует доменный input в ресурс репозитория.
func ToRepoInput(m model.UserInput) resources.UserInput {
	return resources.UserInput{
		Email: m.Email,
		Name:  m.Name,
	}
}
