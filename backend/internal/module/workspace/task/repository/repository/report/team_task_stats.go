package report

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
// Период для done-задач — по completed_at (момент перевода в done).
// Запрос: LEFT JOIN + агрегация вместо коррелированных подзапросов для лучшей производительности при большом числе команд.
func (r *repository) TeamTaskStats(ctx context.Context, tx *sqlx.Tx, since time.Time) ([]*model.TeamTaskStats, error) {
	const query = `
		SELECT t.id AS team_id, t.name AS team_name,
			COUNT(DISTINCT tm.user_id) AS member_count,
			COUNT(tk.id) AS done_tasks_count
		FROM teams t
		LEFT JOIN team_members tm ON tm.team_id = t.id
		LEFT JOIN tasks tk ON tk.team_id = t.id AND tk.status = 'done' AND tk.deleted_at IS NULL AND tk.completed_at >= ?
		WHERE t.deleted_at IS NULL
		GROUP BY t.id, t.name
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
		stats, err := converter.ToTeamTaskStats(rows[i])
		if err != nil {
			return nil, fmt.Errorf("convert row %d: %w", i, err)
		}
		out = append(out, &stats)
	}
	return out, nil
}
