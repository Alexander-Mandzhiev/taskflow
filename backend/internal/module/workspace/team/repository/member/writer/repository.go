package writer

import (
	memberDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/member"
	"github.com/jmoiron/sqlx"
)

var _ memberDef.MemberWriterRepository = (*repository)(nil)

type repository struct {
	writePool *sqlx.DB
}

// NewRepository создаёт репозиторий записи участников. db — тот же *sqlx.DB, что передаётся в TxManager.
func NewRepository(db *sqlx.DB) *repository {
	return &repository{writePool: db}
}
