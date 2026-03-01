package cache

import (
	"context"
	"time"
)

type RedisClient interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error

	// Get возвращает значение по ключу. При промахе (ключ отсутствует) — (nil, nil), не ошибка.
	Get(ctx context.Context, key string) ([]byte, error)

	// MGet возвращает значения по ключам одним запросом. Порядок результатов соответствует порядку keys.

	// При промахе на месте элемента — nil (не ошибка).

	MGet(ctx context.Context, keys ...string) ([][]byte, error)

	Del(ctx context.Context, key string) error

	Ping(ctx context.Context) error

	// Hash operations

	HSet(ctx context.Context, key string, values map[string]interface{}) error

	HGet(ctx context.Context, key, field string) (string, error)

	HGetAll(ctx context.Context, key string) (map[string]string, error)

	HDel(ctx context.Context, key string, fields ...string) error

	Expire(ctx context.Context, key string, ttl time.Duration) error
}
