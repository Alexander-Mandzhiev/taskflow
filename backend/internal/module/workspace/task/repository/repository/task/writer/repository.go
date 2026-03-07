package writer

import (
	"github.com/jmoiron/sqlx"

	taskDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/task"
)

var _ taskDef.TaskWriterRepository = (*repository)(nil)

type repository struct {
	writePool *sqlx.DB
}

// NewRepository создаёт репозиторий записи задач. db — тот же *sqlx.DB, что передаётся в TxManager.
func NewRepository(db *sqlx.DB) *repository {
	return &repository{writePool: db}
}
