package writer

import (
	"github.com/jmoiron/sqlx"

	historyDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/repository/history"
)

var _ historyDef.TaskHistoryWriterRepository = (*repository)(nil)

type repository struct {
	writePool *sqlx.DB
}

// NewRepository создаёт репозиторий записи истории задач.
func NewRepository(db *sqlx.DB) *repository {
	return &repository{writePool: db}
}
