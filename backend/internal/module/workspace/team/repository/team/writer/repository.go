package writer

import (
	"github.com/jmoiron/sqlx"

	teamDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/team"
)

var _ teamDef.TeamWriterRepository = (*repository)(nil)

type repository struct {
	writePool *sqlx.DB
}

// NewRepository создаёт репозиторий записи команд. db — тот же *sqlx.DB, что передаётся в TxManager.
func NewRepository(db *sqlx.DB) *repository {
	return &repository{writePool: db}
}
