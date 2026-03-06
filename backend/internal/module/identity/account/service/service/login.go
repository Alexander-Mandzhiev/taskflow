package service

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	accountmodel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/model"
	usermodel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/jwt"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/useragent"
)

// Login проверяет учётные данные и создаёт сессию в кеше. При неверных данных — accountmodel.ErrInvalidCredentials.
func (s *accountService) Login(ctx context.Context, input accountmodel.LoginInput) (accessToken, refreshToken string, err error) {
	user, err := s.userRepo.GetByEmail(ctx, nil, input.Email)
	if err != nil {
		if errors.Is(err, usermodel.ErrUserNotFound) {
			return "", "", accountmodel.ErrInvalidCredentials
		}
		return "", "", err
	}
	// user не nil по контракту UserRepository

	if err := s.hasher.Compare(user.PasswordHash, input.Password); err != nil {
		return "", "", accountmodel.ErrInvalidCredentials
	}

	client := useragent.DeviceTypeFromUserAgent(input.UserAgent)
	userIDStr := user.ID.String()

	refreshToken, jti, err := jwt.GenerateRefreshToken(userIDStr, client, s.refreshSecret, s.refreshTTL)
	if err != nil {
		logger.Error(ctx, "Login: generate refresh token failed", zap.Error(err))
		return "", "", err
	}

	session := &accountmodel.Session{
		UserID:     user.ID,
		CreatedAt:  time.Now(),
		DeviceType: client,
		UserAgent:  input.UserAgent,
		IP:         input.IP,
	}
	if err := s.sessionRepo.Set(ctx, jti, session, s.sessionTTL); err != nil {
		logger.Error(ctx, "Login: set session failed", zap.Error(err))
		return "", "", err
	}

	accessToken, err = jwt.GenerateToken(userIDStr, client, s.accessSecret, s.accessTTL)
	if err != nil {
		logger.Error(ctx, "Login: generate access token failed", zap.Error(err))
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
