package jwt

import (
	"strings"
	"testing"
	"time"
)

const (
	validUserID   = "550e8400-e29b-41d4-a716-446655440000"
	validClientID = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	secret        = "test-secret-key"
)

func TestGenerateToken_Validation(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		clientID  string
		secretKey string
		expiry    time.Duration
		wantErr   string
	}{
		{"empty userID", "", validClientID, secret, time.Hour, "userID cannot be empty"},
		{"empty clientID", validUserID, "", secret, time.Hour, "clientID cannot be empty"},
		{"empty secretKey", validUserID, validClientID, "", time.Hour, "secretKey cannot be empty"},
		{"zero expiry", validUserID, validClientID, secret, 0, "expiryDuration must be positive"},
		{"negative expiry", validUserID, validClientID, secret, -time.Minute, "expiryDuration must be positive"},
		{"invalid userID", "not-uuid", validClientID, secret, time.Hour, "userID must be a valid UUID"},
		{"invalid clientID", validUserID, "not-uuid", secret, time.Hour, "clientID must be a valid UUID"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GenerateToken(tt.userID, tt.clientID, tt.secretKey, tt.expiry)
			if err == nil {
				t.Fatalf("GenerateToken() expected error containing %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("GenerateToken() error = %v, want substring %q", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateToken_Success(t *testing.T) {
	token, err := GenerateToken(validUserID, validClientID, secret, time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken() err = %v", err)
	}
	if token == "" {
		t.Error("GenerateToken() returned empty token")
	}
	// три части: header.payload.signature
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Errorf("token should have 3 parts, got %d", len(parts))
	}
}

func TestGenerateToken_RoundTrip(t *testing.T) {
	token, err := GenerateToken(validUserID, validClientID, secret, time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("ValidateToken(valid token): %v", err)
	}
	if claims.UserID != validUserID {
		t.Errorf("claims.UserID = %q, want %q", claims.UserID, validUserID)
	}
	if claims.ClientID != validClientID {
		t.Errorf("claims.ClientID = %q, want %q", claims.ClientID, validClientID)
	}
}
