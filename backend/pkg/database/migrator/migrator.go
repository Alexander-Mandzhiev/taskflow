package migrator

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

// Migrator применяет SQL-миграции (goose) к БД.
type Migrator struct {
	db            *sql.DB
	migrationsDir string
	dialect       string
	logger        Logger
}

// NewGooseMigrator создаёт мигратор с указанным логгером. Если log == nil, используется NoopLogger.
func NewGooseMigrator(db *sql.DB, migrationsDir, dialect string, log Logger) *Migrator {
	if log == nil {
		log = &logger.NoopLogger{}
	}
	goose.SetLogger(&gooseAdapter{log: log})
	return &Migrator{
		db:            db,
		migrationsDir: migrationsDir,
		dialect:       dialect,
		logger:        log,
	}
}

type gooseAdapter struct {
	log Logger
}

const goosePrefix = "🪿 [goose] "

func (a *gooseAdapter) Print(v ...interface{}) {
	a.log.Info(context.Background(), goosePrefix+fmt.Sprint(v...), zap.String("component", "goose"))
}

func (a *gooseAdapter) Printf(format string, v ...interface{}) {
	a.log.Info(context.Background(), goosePrefix+fmt.Sprintf(format, v...), zap.String("component", "goose"))
}

func (a *gooseAdapter) Fatalf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	a.log.Error(context.Background(), goosePrefix+msg, zap.String("component", "goose"))
	panic(msg)
}

// Up применяет все доступные миграции в порядке возрастания версии.
func (m *Migrator) Up(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if m.db == nil {
		return sql.ErrConnDone
	}

	m.logInfo(ctx, "🔄 [Migrations] Применяем миграции", zap.String("dir", m.migrationsDir), zap.String("dialect", m.dialect))

	if err := goose.SetDialect(m.dialect); err != nil {
		m.logError(ctx, "❌ [Migrations] Не удалось установить диалект goose", err)
		return err
	}

	if err := goose.Up(m.db, m.migrationsDir); err != nil {
		m.logError(ctx, "❌ [Migrations] Не удалось применить миграции", err)
		return err
	}

	m.logInfo(ctx, "✅ [Migrations] Все миграции успешно применены", zap.String("dir", m.migrationsDir))
	return nil
}

func (m *Migrator) logInfo(ctx context.Context, msg string, fields ...zap.Field) {
	if m.logger != nil {
		m.logger.Info(ctx, msg, fields...)
	}
}

func (m *Migrator) logError(ctx context.Context, msg string, err error) {
	if m.logger != nil {
		m.logger.Error(ctx, msg, zap.Error(err))
	}
}
