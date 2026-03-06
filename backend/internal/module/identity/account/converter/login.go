package converter

import (
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/account/v1/dto"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/model"
)

// LoginRequestToDomain конвертирует DTO запроса логина и метаданные запроса в доменную модель.
func LoginRequestToDomain(req dto.LoginRequest, userAgent, ip string) model.LoginInput {
	return model.LoginInput{
		Email:     req.Email,
		Password:  req.Password,
		UserAgent: userAgent,
		IP:        ip,
	}
}
