package invitation

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

// CreateInvitation создаёт запись приглашения в team_invitations (status=pending, token и expires_at заданы вызывающим).
func (r *Adapter) CreateInvitation(ctx context.Context, tx *sqlx.Tx, inv *model.TeamInvitation) error {
	return r.invitationWriter.Create(ctx, tx, inv)
}
