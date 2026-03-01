package cache

import (
	def "mkk/internal/module/identity/account/repository"
	redisclient "mkk/pkg/cache"
)

var _ def.SessionCacheRepository = (*repository)(nil)

type repository struct {
	redis redisclient.RedisClient
}

// NewRepository создаёт кеш-репозиторий сессий (Redis, ключ session:{id}).
func NewRepository(redis redisclient.RedisClient) *repository {
	return &repository{redis: redis}
}
