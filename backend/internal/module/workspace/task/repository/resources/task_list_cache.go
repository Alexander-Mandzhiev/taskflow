package resources

import (
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// TaskListPageCache — одна закешированная страница списка задач (значение в Redis).
// Сериализуется в JSON. Ключ строится по team_id + фильтру + странице (см. TaskListCacheRepository).
type TaskListPageCache struct {
	Items  []model.Task `json:"items"`
	Total  int          `json:"total"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
}
