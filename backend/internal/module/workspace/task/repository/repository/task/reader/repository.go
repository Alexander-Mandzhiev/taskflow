package reader

import (
	"github.com/jmoiron/sqlx"

	taskDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/task"
)

var _ taskDef.TaskReaderRepository = (*repository)(nil)

type repository struct {
	readPool *sqlx.DB
}

// NewRepository создаёт репозиторий чтения задач. db — тот же *sqlx.DB, что передаётся в TxManager.
func NewRepository(db *sqlx.DB) *repository {
	return &repository{readPool: db}
}
