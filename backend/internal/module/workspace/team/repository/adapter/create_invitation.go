package adapter

import (
	"context"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/jmoiron/sqlx"
)

// CreateInvitation создаёт запись приглашения в team_invitations.

func (r *Repository) CreateInvitation(ctx context.Context, tx *sqlx.Tx, inv *model.TeamInvitation) error {
	return r.invitationWriter.Create(ctx, tx, inv)
}
