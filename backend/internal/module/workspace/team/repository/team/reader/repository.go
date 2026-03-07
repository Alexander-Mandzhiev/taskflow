package reader

import (
	teamDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/team"
	"github.com/jmoiron/sqlx"
)

var _ teamDef.TeamReaderRepository = (*repository)(nil)

type repository struct {
	readPool *sqlx.DB
}

// NewRepository создаёт репозиторий чтения команд. db — тот же *sqlx.DB, что передаётся в TxManager.
func NewRepository(db *sqlx.DB) *repository {
	return &repository{readPool: db}
}
