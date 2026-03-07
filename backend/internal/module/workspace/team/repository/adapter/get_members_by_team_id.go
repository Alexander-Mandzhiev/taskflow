package adapter

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

// GetMembersByTeamID возвращает участников команды по team_id.
func (r *Adapter) GetMembersByTeamID(ctx context.Context, tx *sqlx.Tx, teamID uuid.UUID) ([]*model.TeamMember, error) {
	return r.memberReader.GetByTeamID(ctx, tx, teamID)
}
