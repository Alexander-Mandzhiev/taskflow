package adapter

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
)

// GetPendingInvitationByTeamAndEmail возвращает приглашение со статусом pending для (team_id, email).
func (r *Repository) GetPendingInvitationByTeamAndEmail(ctx context.Context, tx *sqlx.Tx, teamID, email string) (*model.TeamInvitation, error) {
	return r.invitationReader.GetPendingByTeamAndEmail(ctx, tx, teamID, email)
}
