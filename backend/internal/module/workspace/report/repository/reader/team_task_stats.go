package reader

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/converter"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/resources"
)

// TeamTaskStats возвращает для каждой команды: название, кол-во участников, кол-во задач в статусе done за период since.
func (r *repository) TeamTaskStats(ctx context.Context, tx *sqlx.Tx, since time.Time) ([]*model.TeamTaskStats, error) {
	const query = `
		SELECT t.id AS team_id, t.name AS team_name,
			(SELECT COUNT(*) FROM team_members tm WHERE tm.team_id = t.id) AS member_count,
			(SELECT COUNT(*) FROM tasks tk WHERE tk.team_id = t.id AND tk.status = 'done' AND tk.deleted_at IS NULL AND tk.updated_at >= ?) AS done_tasks_count
		FROM teams t
		WHERE t.deleted_at IS NULL
		ORDER BY t.name
	`
	var rows []resources.TeamTaskStatsRow
	if tx != nil {
		if err := tx.SelectContext(ctx, &rows, query, since); err != nil {
			return nil, toDomainError(err)
		}
	} else {
		if err := r.readPool.SelectContext(ctx, &rows, query, since); err != nil {
			return nil, toDomainError(err)
		}
	}

	out := make([]*model.TeamTaskStats, 0, len(rows))
	for i := range rows {
		stats, err := converter.ToDomainTeamTaskStats(rows[i])
		if err != nil {
			return nil, fmt.Errorf("convert row %d: %w", i, err)
		}
		out = append(out, &stats)
	}
	return out, nil
}
