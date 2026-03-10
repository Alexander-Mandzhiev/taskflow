package task

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/resources"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

// List возвращает список задач по фильтру (критерии + limit/offset в filter). total — общее количество без LIMIT.
// При наличии listCache: сначала Get, при попадании — возврат; при промахе — reader.List и Set.
func (r *Adapter) List(ctx context.Context, tx *sqlx.Tx, filter model.TaskListFilter) ([]model.Task, int, error) {
	if r.listCache != nil && filter.TeamID != nil {
		cached, err := r.listCache.Get(ctx, *filter.TeamID, filter)
		if err != nil {
			logger.Warn(ctx, "task list cache get failed, falling back to DB", zap.Error(err))
		} else if cached != nil {
			return cached.Items, cached.Total, nil
		}
	}

	items, total, err := r.taskReader.List(ctx, tx, filter)
	if err != nil {
		return nil, 0, err
	}

	if r.listCache != nil && filter.TeamID != nil {
		if setErr := r.listCache.Set(ctx, *filter.TeamID, filter, &resources.TaskListPageCache{
			Items:  items,
			Total:  total,
			Limit:  filter.Limit,
			Offset: filter.Offset,
		}, 5*time.Minute); setErr != nil {
			logger.Warn(ctx, "task list cache set failed", zap.Error(setErr))
		}
	}
	return items, total, nil
}
