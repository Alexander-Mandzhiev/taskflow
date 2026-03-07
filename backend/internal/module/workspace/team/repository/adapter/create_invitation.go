package adapter

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

// CreateInvitation создаёт запись приглашения в team_invitations.

func (r *Adapter) CreateInvitation(ctx context.Context, tx *sqlx.Tx, inv *model.TeamInvitation) error {
	return r.invitationWriter.Create(ctx, tx, inv)
}
