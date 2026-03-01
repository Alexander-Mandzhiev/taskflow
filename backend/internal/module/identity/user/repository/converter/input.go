package converter

import (
	"mkk/internal/module/identity/user/model"
	"mkk/internal/module/identity/user/repository/resources"
)

// ToRepoInput преобразует доменный input в ресурс репозитория.
// Вызывающий (сервисный слой) гарантирует m != nil; проверка на nil выполняется в сервисе до вызова репозитория.
func ToRepoInput(m *model.UserInput) resources.UserInput {
	return resources.UserInput{
		Email: m.Email,
		Name:  m.Name,
	}
}
