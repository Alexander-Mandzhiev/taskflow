package team

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

// Create создаёт запись в teams (created_by = ownerUserID). Добавление owner в team_members — зона ответственности сервиса (AddMember).
func (r *Adapter) Create(ctx context.Context, tx *sqlx.Tx, input model.TeamInput, ownerUserID uuid.UUID) (model.Team, error) {
	return r.teamWriter.Create(ctx, tx, input, ownerUserID)
}
