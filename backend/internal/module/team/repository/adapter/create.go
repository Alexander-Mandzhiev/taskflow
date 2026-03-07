package adapter

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
)

// Create создаёт запись в teams (created_by = ownerUserID). Добавление owner в team_members выполняет сервис через AddMember.
func (r *Repository) Create(ctx context.Context, tx *sqlx.Tx, input *model.TeamInput, ownerUserID string) (*model.Team, error) {
	return r.teamWriter.Create(ctx, tx, input, ownerUserID)
}
