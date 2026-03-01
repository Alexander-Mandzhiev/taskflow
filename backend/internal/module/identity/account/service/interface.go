package service

import (
	"context"

	"github.com/google/uuid"
)

// AccountService — слой сервиса account: регистрация, логин, логаут, проверка сессии.
// Реализация использует user service, хранилище сессий и хеширование паролей.
type AccountService interface {
	Register(ctx context.Context, email, password, name string) error
	// Login: userAgent и ip — опционально, для отображения в списке сессий (безопасность: пользователь может завершить подозрительную сессию).
	Login(ctx context.Context, email, password, userAgent, ip string) (sessionID uuid.UUID, err error)
	Logout(ctx context.Context, sessionID uuid.UUID) error
	Whoami(ctx context.Context, sessionID uuid.UUID) (userID uuid.UUID, err error)
}

