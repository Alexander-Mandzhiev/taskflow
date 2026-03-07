package reader

import (
	"github.com/jmoiron/sqlx"

	reportDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/report/repository"
)

var _ reportDef.ReportRepository = (*repository)(nil)

type repository struct {
	readPool *sqlx.DB
}

// NewRepository создаёт репозиторий отчётов. db — тот же *sqlx.DB, что и для задач/команд.
func NewRepository(db *sqlx.DB) *repository {
	return &repository{readPool: db}
}
