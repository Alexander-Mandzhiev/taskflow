package jwt

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestValidateToken_EmptySecret(t *testing.T) {
	_, err := ValidateToken("any.token.here", "")
	if err == nil {
		t.Fatal("ValidateToken(_, \"\") expected error")
	}
	if !strings.Contains(err.Error(), "secret key cannot be empty") {
		t.Errorf("ValidateToken error = %v", err)
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	_, err := ValidateToken("not-a-valid-jwt", secret)
	if err == nil {
		t.Fatal("ValidateToken(invalid) expected error")
	}
	if !errors.Is(err, ErrInvalidToken) {
		t.Errorf("ValidateToken(invalid) err = %v, want ErrInvalidToken", err)
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	token, err := GenerateToken(validUserID, validClient, secret, time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	_, err = ValidateToken(token, "wrong-secret")
	if err == nil {
		t.Fatal("ValidateToken(wrong secret) expected error")
	}
	if !errors.Is(err, ErrInvalidToken) {
		t.Errorf("ValidateToken(wrong secret) err = %v, want ErrInvalidToken", err)
	}
}

func TestValidateToken_Expired(t *testing.T) {
	// токен с истёкшим сроком (ExpiresAt в прошлом) — стабильно в CI, без ожидания по таймауту
	now := time.Now()
	expiredAt := now.Add(-time.Hour)
	claims := &Claims{
		Client: validClient,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   validUserID,
			Audience:  jwt.ClaimStrings{validClient},
			ExpiresAt: jwt.NewNumericDate(expiredAt),
			IssuedAt:  jwt.NewNumericDate(now.Add(-2 * time.Hour)),
			Issuer:    "taskflow",
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tok.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign expired token: %v", err)
	}

	_, err = ValidateToken(token, secret)
	if err == nil {
		t.Fatal("ValidateToken(expired) expected error")
	}
	if !errors.Is(err, ErrExpiredToken) {
		t.Errorf("ValidateToken(expired) err = %v, want ErrExpiredToken", err)
	}
}

func TestValidateToken_Success(t *testing.T) {
	token, err := GenerateToken(validUserID, validClient, secret, time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}
	if claims.Subject != validUserID || claims.Client != validClient {
		t.Errorf("claims = Subject %q Client %q, want %q %q", claims.Subject, claims.Client, validUserID, validClient)
	}
}
