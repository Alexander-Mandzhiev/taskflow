package adapter

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
)

// ListByUserID возвращает команды, где пользователь участник, с его ролью в каждой.
func (r *Repository) ListByUserID(ctx context.Context, tx *sqlx.Tx, userID string) ([]*model.TeamWithRole, error) {
	return r.teamReader.ListByUserID(ctx, tx, userID)
}
