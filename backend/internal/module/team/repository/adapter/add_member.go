package adapter

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
)

// AddMember добавляет пользователя в команду с указанной ролью.
func (r *Repository) AddMember(ctx context.Context, tx *sqlx.Tx, teamID, userID, role string) (*model.TeamMember, error) {
	return r.memberWriter.AddMember(ctx, tx, teamID, userID, role)
}
