package adapter

import (
	"context"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// ListByUserID возвращает команды, где пользователь участник, с его ролью в каждой.
func (r *Repository) ListByUserID(ctx context.Context, tx *sqlx.Tx, userID uuid.UUID) ([]*model.TeamWithRole, error) {
	return r.teamReader.ListByUserID(ctx, tx, userID.String())
}
