package report

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// TopCreatorsByTeam возвращает топ-N по созданным задачам в каждой команде за период.
func (r *Adapter) TopCreatorsByTeam(ctx context.Context, tx *sqlx.Tx, since time.Time, limit int) ([]*model.TeamTopCreator, error) {
	return r.reader.TopCreatorsByTeam(ctx, tx, since, limit)
}
