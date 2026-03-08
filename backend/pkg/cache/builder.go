package cache

import (
	"github.com/redis/go-redis/v9"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

// BuildClient создаёт Redis клиент (один узел). Для отключения кеша не передавайте WithAddr или передайте WithAddr("").
// log — реализация Logger (например logger.Logger() или &logger.NoopLogger{}).
func BuildClient(log Logger, tracerName string, options ...Option) (RedisClient, error) {
	cfg := defaultConfig()
	for _, opt := range options {
		opt(cfg)
	}
	if cfg.addr == "" {
		return nil, nil
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:            cfg.addr,
		Password:        cfg.password,
		PoolSize:        cfg.poolSize,
		MinIdleConns:    cfg.minIdleConns,
		PoolTimeout:     cfg.poolTimeout,
		ConnMaxIdleTime: cfg.connMaxIdleTime,
		DialTimeout:     cfg.dialTimeout,
		ReadTimeout:     cfg.readTimeout,
		WriteTimeout:    cfg.writeTimeout,
	})
	if log == nil {
		log = &logger.NoopLogger{}
	}
	scanBatchSize := cfg.scanBatchSize
	if scanBatchSize <= 0 {
		scanBatchSize = 100
	}
	return newClient(&redisAdapter{rdb}, log, cfg.dialTimeout, tracerName, scanBatchSize), nil
}
