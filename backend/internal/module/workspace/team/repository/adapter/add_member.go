package adapter

import (
	"context"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// AddMember добавляет пользователя в команду с указанной ролью.
func (r *Repository) AddMember(ctx context.Context, tx *sqlx.Tx, teamID, userID uuid.UUID, role string) (*model.TeamMember, error) {
	return r.memberWriter.AddMember(ctx, tx, teamID.String(), userID.String(), role)
}
