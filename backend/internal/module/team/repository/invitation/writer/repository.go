package writer

import (
	"github.com/jmoiron/sqlx"

	invitationDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/repository/invitation"
)

var _ invitationDef.InvitationWriterRepository = (*repository)(nil)

type repository struct {
	writePool *sqlx.DB
}

// NewRepository создаёт репозиторий записи приглашений.
func NewRepository(db *sqlx.DB) *repository {
	return &repository{writePool: db}
}
