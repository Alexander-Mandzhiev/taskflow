package service

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	accountmodel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/jwt"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

// Logout завершает сессию: при пустом refreshToken — ничего не делает; иначе валидирует токен, извлекает jti и удаляет сессию из кеша.
func (s *accountService) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}

	claims, err := jwt.ValidateToken(refreshToken, s.refreshSecret)
	if err != nil {
		logger.Debug(ctx, "Logout: invalid refresh token", zap.Error(err))
		return accountmodel.ErrInvalidRefreshToken
	}
	if claims.ID == "" {
		return accountmodel.ErrInvalidRefreshToken
	}

	jti, err := uuid.Parse(claims.ID)
	if err != nil {
		return accountmodel.ErrInvalidRefreshToken
	}

	if err := s.sessionRepo.Delete(ctx, jti); err != nil {
		logger.Error(ctx, "Logout: delete session failed", zap.Error(err))
		return err
	}
	return nil
}
