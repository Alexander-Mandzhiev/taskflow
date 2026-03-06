package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func validateTokenParams(userID, client, secretKey string, expiryDuration time.Duration) error {
	if userID == "" {
		return errors.New("userID cannot be empty")
	}
	if client == "" {
		return errors.New("client cannot be empty")
	}
	if secretKey == "" {
		return errors.New("secretKey cannot be empty")
	}
	if expiryDuration <= 0 {
		return errors.New("expiryDuration must be positive")
	}
	if err := validateUUID(userID); err != nil {
		return errors.New("userID must be a valid UUID")
	}
	return nil
}

// signToken создаёт claims, подписывает JWT и возвращает токен. Если withJTI == true, генерирует jti и записывает в claims (для refresh).
func signToken(userID, client, secretKey string, expiryDuration time.Duration, withJTI bool) (token string, jti uuid.UUID, err error) {
	if err := validateTokenParams(userID, client, secretKey, expiryDuration); err != nil {
		return "", uuid.Nil, err
	}

	now := time.Now()
	exp := now.Add(expiryDuration)

	reg := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(exp),
		IssuedAt:  jwt.NewNumericDate(now),
		Issuer:    "taskflow",
		Subject:   userID,
		Audience:  jwt.ClaimStrings{client},
	}
	if withJTI {
		jti = uuid.New()
		reg.ID = jti.String()
	}

	claims := &Claims{Client: client, RegisteredClaims: reg}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tok.SignedString([]byte(secretKey))
	if err != nil {
		return "", uuid.Nil, err
	}
	return token, jti, nil
}

// GenerateToken создаёт access JWT (без jti).
func GenerateToken(userID, client, secretKey string, expiryDuration time.Duration) (string, error) {
	token, _, err := signToken(userID, client, secretKey, expiryDuration, false)
	return token, err
}

// GenerateRefreshToken создаёт refresh JWT с jti в claims и возвращает jti для ключа сессии в Redis.
func GenerateRefreshToken(userID, client, secretKey string, expiryDuration time.Duration) (string, uuid.UUID, error) {
	return signToken(userID, client, secretKey, expiryDuration, true)
}
