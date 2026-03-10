package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/resources"
)

// Get возвращает закешированную страницу. При промахе или ошибке парсинга — (nil, nil).
func (r *Repository) Get(ctx context.Context, teamID uuid.UUID, filter model.TaskListFilter) (*resources.TaskListPageCache, error) {
	key := Key(teamID, filter)
	data, err := r.redis.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("task list cache get: %w", err)
	}
	if data == nil {
		return nil, nil
	}
	var out resources.TaskListPageCache
	if err := json.Unmarshal(data, &out); err != nil {
		_ = r.redis.Del(ctx, key)
		return nil, nil
	}
	return &out, nil
}
