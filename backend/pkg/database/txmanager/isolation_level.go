package txmanager

import "database/sql"

// isolationLevelToString преобразует уровень изоляции в строку для логов и трейсинга.
func isolationLevelToString(level sql.IsolationLevel) string {
	switch level {
	case sql.LevelDefault:
		return "default"
	case sql.LevelReadUncommitted:
		return "read uncommitted"
	case sql.LevelReadCommitted:
		return "read committed"
	case sql.LevelRepeatableRead:
		return "repeatable read"
	case sql.LevelSerializable:
		return "serializable"
	default:
		return "unknown"
	}
}
