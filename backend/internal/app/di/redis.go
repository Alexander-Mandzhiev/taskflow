package di

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"mkk/pkg/cache"
	"mkk/pkg/logger"
)

const defaultRedisTimeout = 3 * time.Second

// RedisClient возвращает обёртку над Redis с трейсингом и логированием. При первом вызове создаёт клиент и регистрирует закрытие в closer.
func (d *Container) RedisClient(ctx context.Context) (cache.RedisClient, error) {
	if d.redisClient != nil {
		return d.redisClient, nil
	}
	if err := d.requireCloser(); err != nil {
		return nil, err
	}

	redisCfg := d.cfg.Redis()
	addr := redisCfg.Addr()
	if addr == "" {
		return nil, fmt.Errorf("redis addr is empty")
	}

	timeout := redisCfg.Timeout()
	if timeout == 0 {
		timeout = defaultRedisTimeout
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisCfg.Password(),
	})

	client := cache.NewClient(rdb, logger.Logger(), timeout, "cache.redis")

	if err := client.Ping(ctx); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	d.cl.AddNamed("Redis client", func(ctx context.Context) error {
		logger.Info(ctx, "Закрытие Redis клиента")
		return rdb.Close()
	})

	logger.Info(ctx, "Redis клиент создан")
	d.redisClient = client
	d.redisCmdable = rdb
	return d.redisClient, nil
}

// RedisCmdable возвращает go-redis Cmdable (для создания модульных клиентов с разными трейсерами).
func (d *Container) RedisCmdable(ctx context.Context) (redis.Cmdable, error) {
	if d.redisCmdable == nil {
		if _, err := d.RedisClient(ctx); err != nil {
			return nil, err
		}
	}
	return d.redisCmdable, nil
}
