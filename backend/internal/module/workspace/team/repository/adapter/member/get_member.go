package member

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

// GetMember возвращает участника по (team_id, user_id). При отсутствии — (nil, model.ErrMemberNotFound).
func (r *Adapter) GetMember(ctx context.Context, tx *sqlx.Tx, teamID, userID uuid.UUID) (*model.TeamMember, error) {
	return r.memberReader.GetMember(ctx, tx, teamID, userID)
}
