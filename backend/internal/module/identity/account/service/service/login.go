package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	accountmodel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/model"
	usermodel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/useragent"
)

// Login проверяет email/пароль и создаёт сессию в кеше. При неверных данных — accountmodel.ErrInvalidCredentials
// (единое сообщение, без раскрытия «пользователь не найден» vs «неверный пароль»).
// userAgent и ip сохраняются в сессии для списка сессий (пользователь может завершить подозрительную сессию).
func (s *accountService) Login(ctx context.Context, email, password, userAgent, ip string) (sessionID uuid.UUID, err error) {
	user, err := s.userRepo.GetByEmail(ctx, nil, email)
	if err != nil {
		if errors.Is(err, usermodel.ErrUserNotFound) {
			return uuid.Nil, accountmodel.ErrInvalidCredentials
		}
		return uuid.Nil, err
	}
	if user == nil {
		return uuid.Nil, accountmodel.ErrInvalidCredentials
	}

	if err := s.hasher.Compare(user.PasswordHash, password); err != nil {
		return uuid.Nil, accountmodel.ErrInvalidCredentials
	}

	sessionID = uuid.New()
	session := &accountmodel.Session{
		UserID:     user.ID,
		CreatedAt:  time.Now(),
		DeviceType: useragent.DeviceTypeFromUserAgent(userAgent),
		UserAgent:  userAgent,
		IP:         ip,
	}
	if err := s.sessionRepo.Set(ctx, sessionID, session, s.sessionTTL); err != nil {
		logger.Error(ctx, "Login: set session failed", zap.Error(err))
		return uuid.Nil, err
	}

	return sessionID, nil
}
