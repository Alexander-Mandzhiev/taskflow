package writer

import (
	"github.com/jmoiron/sqlx"

	memberDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/member"
)

var _ memberDef.MemberWriterRepository = (*repository)(nil)

type repository struct {
	writePool *sqlx.DB
}

// NewRepository создаёт репозиторий записи участников. db — тот же *sqlx.DB, что передаётся в TxManager.
func NewRepository(db *sqlx.DB) *repository {
	return &repository{writePool: db}
}
