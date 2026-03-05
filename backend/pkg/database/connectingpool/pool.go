package connectingpool

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database"
)

// Pool — пул соединений к БД на основе database/sql, с обёрткой sqlx для маппинга и BeginTxx.
// Единственное место создания *sqlx.DB: при bootstrap передавайте SqlxDB() в TxManager и репозитории.
type Pool struct {
	db *sqlx.DB
}

// New открывает пул соединений и проверяет доступность БД.
// driverName — имя драйвера (например "mysql", "postgres"); используется и для sqlx.
// dsn — строка подключения. Драйвер должен быть зарегистрирован через import _ "driver/package".
func New(ctx context.Context, driverName, dsn string, options ...Option) (*Pool, error) {
	sqlDB, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	cfg := defaultConfig()
	for _, opt := range options {
		opt(cfg)
	}
	sqlDB.SetMaxOpenConns(cfg.maxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.maxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.connMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.connMaxIdleTime)
	if err := sqlDB.PingContext(ctx); err != nil {
		if closeErr := sqlDB.Close(); closeErr != nil {
			return nil, fmt.Errorf("ping db: %w (close: %w)", err, closeErr)
		}
		return nil, fmt.Errorf("ping db: %w", err)
	}
	sqlxDB := sqlx.NewDb(sqlDB, driverName)
	return &Pool{db: sqlxDB}, nil
}

// SqlxDB возвращает *sqlx.DB для TxManager и репозиториев (reader, writer). Основной способ доступа к БД.
func (p *Pool) SqlxDB() *sqlx.DB {
	return p.db
}

// DB возвращает *sql.DB (под капотом того же пула). Для совместимости с кодом, ожидающим только database/sql.
func (p *Pool) DB() *sql.DB {
	return p.db.DB
}

// Close закрывает пул соединений.
func (p *Pool) Close() error {
	return p.db.Close()
}

// Ping проверяет доступность БД.
func (p *Pool) Ping(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

// compile-time check
var _ database.DBProvider = (*Pool)(nil)
