package report

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// TeamTaskStats возвращает для каждой команды: название, кол-во участников, кол-во задач done за период.
func (r *Adapter) TeamTaskStats(ctx context.Context, tx *sqlx.Tx, since time.Time) ([]model.TeamTaskStats, error) {
	return r.reader.TeamTaskStats(ctx, tx, since)
}
