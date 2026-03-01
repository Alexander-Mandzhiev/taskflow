package metadata

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type ctxKey string

const (
	sessionIDKey ctxKey = "session_id"
	userIDKey    ctxKey = "user_id"
)

var ErrNotFound = errors.New("metadata: value not found in context")

// SessionID возвращает session_id из контекста.
func SessionID(ctx context.Context) (uuid.UUID, error) {
	v, ok := ctx.Value(sessionIDKey).(uuid.UUID)
	if !ok || v == uuid.Nil {
		return uuid.Nil, ErrNotFound
	}
	return v, nil
}

// SetSessionID записывает session_id в контекст (id — строка UUID, как из cookie).
func SetSessionID(ctx context.Context, id string) context.Context {
	u, _ := uuid.Parse(id)
	return context.WithValue(ctx, sessionIDKey, u)
}

// UserID возвращает user_id из контекста.
func UserID(ctx context.Context) (uuid.UUID, error) {
	v, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok || v == uuid.Nil {
		return uuid.Nil, ErrNotFound
	}
	return v, nil
}

// SetUserID записывает user_id в контекст. id — строка (UUID), как в auth middleware.
func SetUserID(ctx context.Context, id string) context.Context {
	u, _ := uuid.Parse(id)
	return context.WithValue(ctx, userIDKey, u)
}

// SetUserIDUUID записывает user_id (uuid.UUID) в контекст.
func SetUserIDUUID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}
