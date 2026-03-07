package adapter

import (
	"context"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// GetByID возвращает только команду. При отсутствии — model.ErrTeamNotFound.
func (r *Repository) GetByID(ctx context.Context, tx *sqlx.Tx, teamID uuid.UUID) (*model.Team, error) {
	team, err := r.teamReader.GetByID(ctx, tx, teamID.String())
	if err != nil {
		return nil, err
	}

	return team, nil
}
