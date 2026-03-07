package reader

import (
	"github.com/jmoiron/sqlx"

	historyDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/repository/history"
)

var _ historyDef.TaskHistoryReaderRepository = (*repository)(nil)

type repository struct {
	readPool *sqlx.DB
}

// NewRepository создаёт репозиторий чтения истории задач.
func NewRepository(db *sqlx.DB) *repository {
	return &repository{readPool: db}
}
