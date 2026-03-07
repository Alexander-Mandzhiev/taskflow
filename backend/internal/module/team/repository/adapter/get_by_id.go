package adapter

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
)

// GetByID возвращает команду с участниками. Оба чтения в одном tx при tx != nil — консистентный снимок.
func (r *Repository) GetByID(ctx context.Context, tx *sqlx.Tx, teamID string) (*model.TeamWithMembers, error) {
	team, err := r.teamReader.GetByID(ctx, tx, teamID)
	if err != nil {
		return nil, err
	}

	members, err := r.memberReader.GetByTeamID(ctx, tx, teamID)
	if err != nil {
		return nil, err
	}

	return &model.TeamWithMembers{Team: *team, Members: members}, nil
}
