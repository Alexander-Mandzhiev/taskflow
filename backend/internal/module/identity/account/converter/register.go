package converter

import (
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/account/v1/dto"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/model"
)

// RegisterRequestToDomain конвертирует DTO запроса регистрации в доменную модель.
func RegisterRequestToDomain(req dto.RegisterRequest) model.RegisterInput {
	return model.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}
}
