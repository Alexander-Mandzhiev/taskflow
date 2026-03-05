package writer

import (
	"github.com/jmoiron/sqlx"

	def "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository"
)

var _ def.UserWriterRepository = (*repository)(nil)

type repository struct {
	writePool *sqlx.DB
}

// NewRepository создаёт репозиторий записи. db — тот же *sqlx.DB, что передаётся в TxManager (например pool.SqlxDB()).
func NewRepository(db *sqlx.DB) *repository {
	return &repository{writePool: db}
}
