package team

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

// ListByUserID возвращает команды, где пользователь участник, с его ролью в каждой (GET /api/v1/teams).
func (r *Adapter) ListByUserID(ctx context.Context, tx *sqlx.Tx, userID uuid.UUID) ([]*model.TeamWithRole, error) {
	return r.teamReader.ListByUserID(ctx, tx, userID)
}
