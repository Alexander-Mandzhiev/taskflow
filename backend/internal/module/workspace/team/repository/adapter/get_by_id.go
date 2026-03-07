package adapter

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

// GetByID возвращает только команду. При отсутствии — model.ErrTeamNotFound.
func (r *Repository) GetByID(ctx context.Context, tx *sqlx.Tx, teamID uuid.UUID) (*model.Team, error) {
	return r.teamReader.GetByID(ctx, tx, teamID)
}
