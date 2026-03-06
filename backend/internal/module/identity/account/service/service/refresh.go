package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"

	accountmodel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/jwt"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/useragent"
)

// Refresh по валидному refresh-токену выдаёт новый access-токен и userID. Проверяет сессию в Redis по jti и мягко сверяет userAgent/IP.
func (s *accountService) Refresh(ctx context.Context, refreshToken, userAgent, ip string) (accessToken string, userID uuid.UUID, err error) {
	if refreshToken == "" {
		return "", uuid.Nil, accountmodel.ErrInvalidRefreshToken
	}

	claims, err := jwt.ValidateToken(refreshToken, s.refreshSecret)
	if err != nil {
		logger.Debug(ctx, "Refresh: invalid refresh token", zap.Error(err))
		return "", uuid.Nil, accountmodel.ErrInvalidRefreshToken
	}
	if claims.Subject == "" || claims.Client == "" || claims.ID == "" {
		return "", uuid.Nil, accountmodel.ErrInvalidRefreshToken
	}

	jti, err := uuid.Parse(claims.ID)
	if err != nil {
		return "", uuid.Nil, accountmodel.ErrInvalidRefreshToken
	}

	session, err := s.sessionRepo.Get(ctx, jti)
	if err != nil {
		if errors.Is(err, accountmodel.ErrSessionNotFound) {
			logger.Debug(ctx, "Refresh: session not found or expired in Redis", zap.String("jti", jti.String()))
			return "", uuid.Nil, accountmodel.ErrInvalidRefreshToken
		}
		logger.Error(ctx, "Refresh: get session failed", zap.String("jti", jti.String()), zap.Error(err))
		return "", uuid.Nil, err
	}

	if !sessionMatchesRequest(session, userAgent, ip) {
		logger.Debug(ctx, "Refresh: session metadata mismatch (userAgent/IP)")
		return "", uuid.Nil, accountmodel.ErrInvalidRefreshToken
	}

	parsedUserID, err := uuid.Parse(claims.Subject)
	if err != nil {
		logger.Debug(ctx, "Refresh: invalid subject in claims (not a valid UUID)", zap.String("subject", claims.Subject), zap.Error(err))
		return "", uuid.Nil, accountmodel.ErrInvalidRefreshToken
	}

	accessToken, err = jwt.GenerateToken(claims.Subject, claims.Client, s.accessSecret, s.accessTTL)
	if err != nil {
		logger.Error(ctx, "Refresh: generate access token failed", zap.Error(err))
		return "", uuid.Nil, err
	}

	return accessToken, parsedUserID, nil
}

// sessionMatchesRequest — мягкая сверка: тип устройства по User-Agent и при необходимости IP.
func sessionMatchesRequest(session *accountmodel.Session, userAgent, ip string) bool {
	if session.UserAgent != "" && userAgent != "" {
		if useragent.DeviceTypeFromUserAgent(userAgent) != session.DeviceType {
			return false
		}
	}
	if session.IP != "" && ip != "" {
		if session.IP != ip {
			return false
		}
	}
	return true
}
