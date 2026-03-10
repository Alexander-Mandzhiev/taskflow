package list

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/resources"
)

// Set сохраняет страницу в кеш с TTL.
func (r *Repository) Set(ctx context.Context, teamID uuid.UUID, filter model.TaskListFilter, data *resources.TaskListPageCache, ttl time.Duration) error {
	if data == nil {
		return nil
	}
	if ttl <= 0 {
		ttl = TTL
	}
	key := Key(teamID, filter)
	raw, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("task list cache marshal: %w", err)
	}
	if err := r.redis.Set(ctx, key, raw, ttl); err != nil {
		return fmt.Errorf("task list cache set: %w", err)
	}
	return nil
}
