package writer

import (
	invitationDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/invitation"
	"github.com/jmoiron/sqlx"
)

var _ invitationDef.InvitationWriterRepository = (*repository)(nil)

type repository struct {
	writePool *sqlx.DB
}

// NewRepository создаёт репозиторий записи приглашений.
func NewRepository(db *sqlx.DB) *repository {
	return &repository{writePool: db}
}
