package metadata

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/ctxkey"
)

var ErrNotFound = errors.New("metadata: value not found in context")

// SessionID возвращает session_id из контекста.
func SessionID(ctx context.Context) (uuid.UUID, error) {
	v, ok := ctx.Value(ctxkey.SessionID).(uuid.UUID)
	if !ok || v == uuid.Nil {
		return uuid.Nil, ErrNotFound
	}
	return v, nil
}

// SetSessionIDUUID записывает session_id (uuid.UUID) в контекст.
func SetSessionIDUUID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, ctxkey.SessionID, id)
}

// UserID возвращает user_id из контекста.
func UserID(ctx context.Context) (uuid.UUID, error) {
	v, ok := ctx.Value(ctxkey.UserID).(uuid.UUID)
	if !ok || v == uuid.Nil {
		return uuid.Nil, ErrNotFound
	}
	return v, nil
}

// SetUserID записывает user_id в контекст. id — строка (UUID), как в auth middleware.
func SetUserID(ctx context.Context, id string) context.Context {
	u, _ := uuid.Parse(id)
	return context.WithValue(ctx, ctxkey.UserID, u)
}

// SetUserIDUUID записывает user_id (uuid.UUID) в контекст.
func SetUserIDUUID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, ctxkey.UserID, id)
}
