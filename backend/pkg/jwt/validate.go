package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ValidateToken валидирует JWT токен и возвращает claims.
// Проверяет подпись, срок действия и формат claims.
func ValidateToken(tokenString, secretKey string) (*Claims, error) {
	if secretKey == "" {
		return nil, errors.New("secret key cannot be empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidToken
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// Валидация обязательных полей
	if claims.UserID == "" {
		return nil, ErrInvalidToken
	}
	if claims.ClientID == "" {
		return nil, ErrInvalidToken
	}

	// Валидация формата UUID
	if err := validateUUID(claims.UserID); err != nil {
		return nil, ErrInvalidToken
	}
	if err := validateUUID(claims.ClientID); err != nil {
		return nil, ErrInvalidToken
	}

	// NOTE: Проверка срока действия уже выполняется jwt.ParseWithClaims
	// Дополнительная проверка ниже оставлена как защита на случай изменений в библиотеке
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, ErrExpiredToken
	}

	return claims, nil
}

// validateUUID вспомогательная функция для валидации UUID
func validateUUID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return err
	}
	return nil
}
