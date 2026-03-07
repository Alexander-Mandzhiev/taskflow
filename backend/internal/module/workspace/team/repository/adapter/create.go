package adapter

import (
	"context"

	model2 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Create создаёт запись в teams (created_by = ownerUserID). Добавление owner в team_members выполняет сервис через AddMember.
func (r *Repository) Create(ctx context.Context, tx *sqlx.Tx, input *model2.TeamInput, ownerUserID uuid.UUID) (*model2.Team, error) {
	return r.teamWriter.Create(ctx, tx, input, ownerUserID.String())
}
