package model

import "errors"

// Ошибки домена account (сессии). Сервис и API проверяют через errors.Is.
var (
	// ErrSessionNotFound — сессия не найдена или истекла (кеш/хранилище сессий).
	ErrSessionNotFound = errors.New("session not found")
	// ErrInvalidRefreshToken — refresh-токен невалиден или истёк (подпись, срок, формат).
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
	// ErrInvalidCredentials — неверный email или пароль (единое сообщение, не раскрываем причину).
	ErrInvalidCredentials = errors.New("invalid email or password")
)
