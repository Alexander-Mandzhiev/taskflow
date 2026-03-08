package list

import (
	"time"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository"
	cachepkg "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/cache"
)

// TTL кеша списка задач (готовые страницы).
const TTL = 5 * time.Minute

var _ repository.TaskListCacheRepository = (*Repository)(nil)

// Repository реализует TaskListCacheRepository через Redis (Get/Set/DelByPrefix).
type Repository struct {
	redis cachepkg.RedisClient
}

// NewRepository создаёт кеш-репозиторий списка задач.
func NewRepository(redis cachepkg.RedisClient) *Repository {
	return &Repository{redis: redis}
}
