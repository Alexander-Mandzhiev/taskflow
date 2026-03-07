package adapter

import (
	"context"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// GetPendingInvitationByTeamAndEmail возвращает приглашение со статусом pending для (team_id, email).
func (r *Repository) GetPendingInvitationByTeamAndEmail(ctx context.Context, tx *sqlx.Tx, teamID uuid.UUID, email string) (*model.TeamInvitation, error) {
	return r.invitationReader.GetPendingByTeamAndEmail(ctx, tx, teamID.String(), email)
}
