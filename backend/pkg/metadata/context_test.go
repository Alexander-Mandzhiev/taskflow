package metadata

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/ctxkey"
)

func TestSessionKey(t *testing.T) {
	t.Run("empty context returns ErrNotFound", func(t *testing.T) {
		ctx := context.Background()
		id, err := SessionKey(ctx)
		if id != uuid.Nil {
			t.Errorf("SessionKey(empty ctx) id = %v, want Nil", id)
		}
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("SessionKey(empty ctx) err = %v, want ErrNotFound", err)
		}
	})

	t.Run("context with nil UUID returns ErrNotFound", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), ctxkey.SessionKey, uuid.Nil)
		id, err := SessionKey(ctx)
		if id != uuid.Nil {
			t.Errorf("SessionKey(nil UUID) id = %v, want Nil", id)
		}
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("SessionKey(nil UUID) err = %v, want ErrNotFound", err)
		}
	})

	t.Run("SetSessionKey then SessionKey returns value", func(t *testing.T) {
		want := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
		ctx := SetSessionKey(context.Background(), want)
		got, err := SessionKey(ctx)
		if err != nil {
			t.Fatalf("SessionKey: %v", err)
		}
		if got != want {
			t.Errorf("SessionKey = %v, want %v", got, want)
		}
	})
}

func TestUserID(t *testing.T) {
	t.Run("empty context returns ErrNotFound", func(t *testing.T) {
		ctx := context.Background()
		id, err := UserID(ctx)
		if id != uuid.Nil {
			t.Errorf("UserID(empty ctx) id = %v, want Nil", id)
		}
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("UserID(empty ctx) err = %v, want ErrNotFound", err)
		}
	})

	t.Run("context with nil UUID returns ErrNotFound", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), ctxkey.UserID, uuid.Nil)
		id, err := UserID(ctx)
		if id != uuid.Nil {
			t.Errorf("UserID(nil UUID) id = %v, want Nil", id)
		}
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("UserID(nil UUID) err = %v, want ErrNotFound", err)
		}
	})

	t.Run("SetUserIDUUID then UserID returns value", func(t *testing.T) {
		want := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		ctx := SetUserIDUUID(context.Background(), want)
		got, err := UserID(ctx)
		if err != nil {
			t.Fatalf("UserID: %v", err)
		}
		if got != want {
			t.Errorf("UserID = %v, want %v", got, want)
		}
	})

	t.Run("SetUserID with valid string then UserID returns value", func(t *testing.T) {
		str := "7c9e6679-7425-40de-944b-e07fc1f90ae7"
		want := uuid.MustParse(str)
		ctx := SetUserID(context.Background(), str)
		got, err := UserID(ctx)
		if err != nil {
			t.Fatalf("UserID: %v", err)
		}
		if got != want {
			t.Errorf("UserID = %v, want %v", got, want)
		}
	})

	t.Run("SetUserID with invalid string then UserID returns ErrNotFound", func(t *testing.T) {
		ctx := SetUserID(context.Background(), "not-a-uuid")
		got, err := UserID(ctx)
		if got != uuid.Nil {
			t.Errorf("UserID(invalid) id = %v, want Nil", got)
		}
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("UserID(invalid) err = %v, want ErrNotFound", err)
		}
	})
}
