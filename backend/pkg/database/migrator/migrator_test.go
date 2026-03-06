package migrator

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

func TestNewGooseMigrator_NilLogger(t *testing.T) {
	m := NewGooseMigrator(nil, "dir", "mysql", nil)
	if m == nil {
		t.Fatal("NewGooseMigrator(nil db, nil log) = nil")
	}
	if m.logger == nil {
		t.Error("logger should be replaced with NoopLogger")
	}
}

func TestMigrator_Up_NilDB(t *testing.T) {
	m := NewGooseMigrator(nil, "migrations", "mysql", &logger.NoopLogger{})
	ctx := context.Background()
	err := m.Up(ctx)
	if !errors.Is(err, sql.ErrConnDone) {
		t.Errorf("Up(nil db) = %v, want sql.ErrConnDone", err)
	}
}

func TestMigrator_Up_CancelledContext(t *testing.T) {
	m := NewGooseMigrator(nil, "migrations", "mysql", &logger.NoopLogger{})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := m.Up(ctx)
	if err == nil {
		t.Fatal("Up(cancelled ctx) expected error")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Up(cancelled ctx) err = %v, want context.Canceled", err)
	}
}
