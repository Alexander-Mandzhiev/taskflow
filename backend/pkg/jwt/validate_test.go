package jwt

import (
	"errors"
	"strings"
	"testing"
	"time"
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
	token, err := GenerateToken(validUserID, validClientID, secret, time.Hour)
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
	// генерируем токен с истечением через 1 мс
	token, err := GenerateToken(validUserID, validClientID, secret, time.Millisecond)
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}
	time.Sleep(2 * time.Millisecond)

	_, err = ValidateToken(token, secret)
	if err == nil {
		t.Fatal("ValidateToken(expired) expected error")
	}
	if !errors.Is(err, ErrExpiredToken) {
		t.Errorf("ValidateToken(expired) err = %v, want ErrExpiredToken", err)
	}
}

func TestValidateToken_Success(t *testing.T) {
	token, err := GenerateToken(validUserID, validClientID, secret, time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}
	if claims.UserID != validUserID || claims.ClientID != validClientID {
		t.Errorf("claims = UserID %q ClientID %q, want %q %q", claims.UserID, claims.ClientID, validUserID, validClientID)
	}
}
