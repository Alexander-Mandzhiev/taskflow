package di

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/cache"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

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

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisCfg.Password(),
	})

	client := cache.NewClient(rdb, logger.Logger(), redisCfg.Timeout(), "cache.redis", 100)

	if err := client.Ping(ctx); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	d.cl.Add(func(ctx context.Context) error {
		err := rdb.Close()
		logger.Info(ctx, "🔴 [Shutdown] Closed Redis client")
		return err
	})

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
