package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken создает новый JWT токен с указанными claims.
// userID и clientID должны быть валидными UUID.
// expiryDuration определяет срок действия токена.
func GenerateToken(userID, clientID, secretKey string, expiryDuration time.Duration) (string, error) {
	// Валидация входных параметров
	if userID == "" {
		return "", errors.New("userID cannot be empty")
	}
	if clientID == "" {
		return "", errors.New("clientID cannot be empty")
	}
	if secretKey == "" {
		return "", errors.New("secretKey cannot be empty")
	}

	// Валидация expiryDuration
	if expiryDuration <= 0 {
		return "", errors.New("expiryDuration must be positive")
	}

	// Валидация формата UUID
	if err := validateUUID(userID); err != nil {
		return "", errors.New("userID must be a valid UUID")
	}
	if err := validateUUID(clientID); err != nil {
		return "", errors.New("clientID must be a valid UUID")
	}

	// Используем один вызов time.Now() для консистентности
	now := time.Now()
	expirationTime := now.Add(expiryDuration)

	claims := &Claims{
		UserID:   userID,
		ClientID: clientID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "mkk",
			Subject:   userID,
			Audience:  jwt.ClaimStrings{clientID},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
