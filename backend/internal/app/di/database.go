package di

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/connectingpool"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/migrator"
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

	d.cl.Add(func(ctx context.Context) error {
		err := pool.Close()
		logger.Info(ctx, "🔌 [Shutdown] Closed MySQL pool")
		return err
	})

	d.dbPool = pool
	d.sqlxDB = pool.SqlxDB()
	return d.sqlxDB, nil
}

// RunMigrations применяет миграции goose. Путь к каталогу — из конфига (app.migrations_dir / MIGRATIONS_DIR).
func (d *Container) RunMigrations(ctx context.Context) error {
	db, err := d.SqlxDB(ctx)
	if err != nil {
		return fmt.Errorf("mysql pool: %w", err)
	}
	dir := d.cfg.App().MigrationsDir()
	if dir == "" {
		return fmt.Errorf("migrations dir is empty (set app.migrations_dir or MIGRATIONS_DIR)")
	}
	m := migrator.NewGooseMigrator(db.DB, dir, "mysql", logger.Logger())
	if err := m.Up(ctx); err != nil {
		return fmt.Errorf("migrations: %w", err)
	}
	return nil
}
