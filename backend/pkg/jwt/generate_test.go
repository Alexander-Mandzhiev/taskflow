package jwt

import (
	"strings"
	"testing"
	"time"
)

const (
	validUserID = "550e8400-e29b-41d4-a716-446655440000"
	validClient = "web"
	secret      = "test-secret-key"
)

func TestGenerateToken_Validation(t *testing.T) {
	tests := []struct {
		name string

		userID string

		clientID string

		secretKey string

		expiry time.Duration

		wantErr string
	}{
		{"empty userID", "", validClient, secret, time.Hour, "userID cannot be empty"},
		{"empty client", validUserID, "", secret, time.Hour, "client cannot be empty"},
		{"empty secretKey", validUserID, validClient, "", time.Hour, "secretKey cannot be empty"},
		{"zero expiry", validUserID, validClient, secret, 0, "expiryDuration must be positive"},
		{"negative expiry", validUserID, validClient, secret, -time.Minute, "expiryDuration must be positive"},
		{"invalid userID", "not-uuid", validClient, secret, time.Hour, "userID must be a valid UUID"},
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
	token, err := GenerateToken(validUserID, validClient, secret, time.Hour)
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
	token, err := GenerateToken(validUserID, validClient, secret, time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("ValidateToken(valid token): %v", err)
	}

	if claims.Subject != validUserID {
		t.Errorf("claims.Subject = %q, want %q", claims.Subject, validUserID)
	}
	if claims.Client != validClient {
		t.Errorf("claims.Client = %q, want %q", claims.Client, validClient)
	}
}
