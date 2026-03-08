package reader

import (
	"github.com/jmoiron/sqlx"

	invitationDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/repository/invitation"
)

var _ invitationDef.InvitationReaderRepository = (*repository)(nil)

type repository struct {
	readPool *sqlx.DB
}

// NewRepository создаёт репозиторий чтения приглашений.
func NewRepository(db *sqlx.DB) *repository {
	return &repository{readPool: db}
}
