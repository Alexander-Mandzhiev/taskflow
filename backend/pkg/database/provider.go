package database

import "database/sql"

// DBProvider — интерфейс для получения *sql.DB.
// Реализуется connectingpool.Pool; позволяет подменять в тестах.
type DBProvider interface {
	DB() *sql.DB
}
