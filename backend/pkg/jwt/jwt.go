package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

// Claims структура claims для JWT токена
// Содержит только идентификаторы пользователя и клиента.
// Права и роли получаются из RBAC для большей безопасности.
type Claims struct {
	UserID   string `json:"user_id"`   // UUID пользователя (обязательный)
	ClientID string `json:"client_id"` // UUID клиента (обязательный)
	jwt.RegisteredClaims
}
