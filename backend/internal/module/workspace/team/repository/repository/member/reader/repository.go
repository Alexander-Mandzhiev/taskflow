package reader

import (
	"github.com/jmoiron/sqlx"

	memberDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/repository/member"
)

var _ memberDef.MemberReaderRepository = (*repository)(nil)

type repository struct {
	readPool *sqlx.DB
}

// NewRepository создаёт репозиторий чтения участников. db — тот же *sqlx.DB, что передаётся в TxManager.
func NewRepository(db *sqlx.DB) *repository {
	return &repository{readPool: db}
}
