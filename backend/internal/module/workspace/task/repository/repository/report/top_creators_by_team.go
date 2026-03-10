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

// TopCreatorsByTeam возвращает топ-N пользователей по количеству созданных задач в каждой команде за период (MySQL 8+ ROW_NUMBER).
func (r *repository) TopCreatorsByTeam(ctx context.Context, tx *sqlx.Tx, since time.Time, limit int) ([]model.TeamTopCreator, error) {
	if limit <= 0 {
		limit = 3
	}
	const query = `
		SELECT team_id, user_id, ` + "`rank`" + `, created_count FROM (
			SELECT team_id, user_id, created_count,
			       ROW_NUMBER() OVER (PARTITION BY team_id ORDER BY created_count DESC) AS ` + "`rank`" + `
			FROM (
				SELECT team_id, created_by AS user_id, COUNT(*) AS created_count
				FROM tasks
				WHERE deleted_at IS NULL AND created_at >= ?
				GROUP BY team_id, created_by
			) AS counted
		) AS ranked
		WHERE ` + "`rank`" + ` <= ?
		ORDER BY team_id, ` + "`rank`" + `
	`
	var rows []resources.TeamTopCreatorRow
	if tx != nil {
		if err := tx.SelectContext(ctx, &rows, query, since, limit); err != nil {
			return nil, toDomainError(err)
		}
	} else {
		if err := r.readPool.SelectContext(ctx, &rows, query, since, limit); err != nil {
			return nil, toDomainError(err)
		}
	}

	out := make([]model.TeamTopCreator, 0, len(rows))
	for i := range rows {
		item, err := converter.ToTeamTopCreator(rows[i])
		if err != nil {
			return nil, fmt.Errorf("convert row %d: %w", i, err)
		}
		out = append(out, item)
	}
	return out, nil
}
