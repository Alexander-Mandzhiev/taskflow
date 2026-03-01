package account

import (
	"context"

	"github.com/google/uuid"
)

// Service — сервисный слой account: регистрация, логин, логаут, проверка сессии.
// Реализация использует user service, хранилище сессий и хеширование паролей.
type Service interface {
	Register(ctx context.Context, email, password, name string) (userID uuid.UUID, err error)
	Login(ctx context.Context, email, password string) (sessionID uuid.UUID, err error)
	Logout(ctx context.Context, sessionID uuid.UUID) error
	Whoami(ctx context.Context, sessionID uuid.UUID) (userID uuid.UUID, err error)
}
