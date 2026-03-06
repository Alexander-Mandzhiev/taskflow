package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

// Claims — расширение стандартных claims (jwt.RegisteredClaims: sub, jti, exp, iss, aud и т.д.).
// Добавляем только Client (тип устройства). Sub = Subject, jti = ID — из пакета.
type Claims struct {
	Client string `json:"client"` // тип устройства: web, mobile, desktop (обязательный)
	jwt.RegisteredClaims
}
