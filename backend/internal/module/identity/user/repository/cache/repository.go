package cache

import (
	def "mkk/internal/module/identity/user/repository"
	redisclient "mkk/pkg/cache"
)

var _ def.UserCacheRepository = (*repository)(nil)

type repository struct {
	redis redisclient.RedisClient
}

// NewRepository создаёт кеш-репозиторий пользователей (Redis).
func NewRepository(redis redisclient.RedisClient) *repository {
	return &repository{redis: redis}
}
