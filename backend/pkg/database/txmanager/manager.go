package txmanager

import (
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Manager — реализация TxManager с трейсингом и post-commit hooks.
// Использует *sqlx.DB для BeginTxx, чтобы по стеку передавать *sqlx.Tx (без обёрток в репо).
type Manager struct {
	db     *sqlx.DB
	tracer trace.Tracer
}

// New создаёт менеджер транзакций для данного *sqlx.DB.
// tracerName — имя трейсера для спанов (например "user.txmanager", "mkk.txmanager").
func New(db *sqlx.DB, tracerName string) *Manager {
	if tracerName == "" {
		tracerName = "txmanager"
	}
	return &Manager{
		db:     db,
		tracer: otel.GetTracerProvider().Tracer(tracerName),
	}
}
