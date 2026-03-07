package adapter

import (
	"context"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// GetMembersByTeamID возвращает участников команды по team_id.
func (r *Repository) GetMembersByTeamID(ctx context.Context, tx *sqlx.Tx, teamID uuid.UUID) ([]*model.TeamMember, error) {
	return r.memberReader.GetByTeamID(ctx, tx, teamID.String())
}
