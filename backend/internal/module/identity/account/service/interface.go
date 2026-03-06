package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/model"
)

// AccountService — слой сервиса account: регистрация, логин, логаут, проверка сессии.
// Реализация использует user repo, хранилище сессий и хеширование паролей.
type AccountService interface {
	Register(ctx context.Context, input model.RegisterInput) error
	// Login возвращает access и refresh токены; при ошибке — пустые строки и err.
	Login(ctx context.Context, input model.LoginInput) (accessToken, refreshToken string, err error)
	// Logout принимает refreshToken (сырая строка из cookie). Пустая строка — ничего не делать; иначе валидация и удаление сессии по jti.
	Logout(ctx context.Context, refreshToken string) error
	Whoami(ctx context.Context, sessionID uuid.UUID) (userID uuid.UUID, err error)
}
