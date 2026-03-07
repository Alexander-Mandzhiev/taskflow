package reader

import (
	invitationDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/invitation"
	"github.com/jmoiron/sqlx"
)

var _ invitationDef.InvitationReaderRepository = (*repository)(nil)

type repository struct {
	readPool *sqlx.DB
}

// NewRepository создаёт репозиторий чтения приглашений.
func NewRepository(db *sqlx.DB) *repository {
	return &repository{readPool: db}
}
