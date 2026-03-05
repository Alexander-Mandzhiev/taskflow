package reader

import (
	"github.com/jmoiron/sqlx"

	def "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository"
)

var _ def.UserReaderRepository = (*repository)(nil)

type repository struct {
	readPool *sqlx.DB
}

// NewRepository создаёт репозиторий чтения. db — тот же *sqlx.DB, что передаётся в TxManager (например pool.SqlxDB()).
func NewRepository(db *sqlx.DB) *repository {
	return &repository{readPool: db}
}
