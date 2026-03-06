package cache

import (
	def "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository/cache"
	redisclient "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/cache"
)

var _ def.UserCacheRepository = (*repository)(nil)

type repository struct {
	redis redisclient.RedisClient
}

// NewRepository создаёт кеш-репозиторий пользователей (Redis).
func NewRepository(redis redisclient.RedisClient) *repository {
	return &repository{redis: redis}
}
