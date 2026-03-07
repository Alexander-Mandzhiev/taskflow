package report

import (
	"github.com/jmoiron/sqlx"
)

var _ ReportReaderRepository = (*repository)(nil)

type repository struct {
	readPool *sqlx.DB
}

// NewRepository создаёт репозиторий отчётов.
func NewRepository(db *sqlx.DB) *repository {
	return &repository{readPool: db}
}
