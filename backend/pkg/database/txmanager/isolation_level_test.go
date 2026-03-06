package txmanager

import (
	"database/sql"
	"testing"
)

func TestIsolationLevelToString(t *testing.T) {
	tests := []struct {
		level sql.IsolationLevel
		want  string
	}{
		{sql.LevelDefault, "default"},
		{sql.LevelReadUncommitted, "read uncommitted"},
		{sql.LevelReadCommitted, "read committed"},
		{sql.LevelRepeatableRead, "repeatable read"},
		{sql.LevelSerializable, "serializable"},
		{sql.IsolationLevel(99), "unknown"},
	}
	for _, tt := range tests {
		got := isolationLevelToString(tt.level)
		if got != tt.want {
			t.Errorf("isolationLevelToString(%v) = %q, want %q", tt.level, got, tt.want)
		}
	}
}
