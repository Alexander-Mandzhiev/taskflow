package di

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/connectingpool"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

// SqlxDB возвращает *sqlx.DB для MySQL. При первом вызове создаёт пул, пингует БД и регистрирует закрытие в closer.
func (d *Container) SqlxDB(ctx context.Context) (*sqlx.DB, error) {
	if d.sqlxDB != nil {
		return d.sqlxDB, nil
	}
	if err := d.requireCloser(); err != nil {
		return nil, err
	}

	mysql := d.cfg.MySQL()
	dsn := mysql.DSN()
	if dsn == "" {
		return nil, fmt.Errorf("mysql dsn is empty")
	}

	opts := []connectingpool.Option{
		connectingpool.WithMaxOpenConns(mysql.MaxOpenConns()),
		connectingpool.WithMaxIdleConns(mysql.MaxIdleConns()),
		connectingpool.WithConnMaxLifetime(mysql.ConnMaxLifetime()),
		connectingpool.WithConnMaxIdleTime(mysql.ConnMaxIdleTime()),
	}
	pool, err := connectingpool.New(ctx, "mysql", dsn, opts...)
	if err != nil {
		return nil, fmt.Errorf("create mysql pool: %w", err)
	}

	d.cl.AddNamed("MySQL pool", func(ctx context.Context) error {
		logger.Info(ctx, "Закрытие MySQL pool")
		return pool.Close()
	})

	logger.Info(ctx, "MySQL pool создан")
	d.dbPool = pool
	d.sqlxDB = pool.SqlxDB()
	return d.sqlxDB, nil
}
