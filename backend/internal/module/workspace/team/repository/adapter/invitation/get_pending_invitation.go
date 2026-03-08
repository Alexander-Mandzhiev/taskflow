package invitation

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

// GetPendingInvitationByTeamAndEmail возвращает приглашение со статусом pending для (team_id, email) или model.ErrInvitationNotFound.
func (r *Adapter) GetPendingInvitationByTeamAndEmail(ctx context.Context, tx *sqlx.Tx, teamID uuid.UUID, email string) (*model.TeamInvitation, error) {
	return r.invitationReader.GetPendingByTeamAndEmail(ctx, tx, teamID, email)
}
