package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/tracing"
)

// internalClient — минимальный набор методов Redis, используемых пакетом. Реализуют *redis.Client и *redis.ClusterClient через адаптер.
type internalClient interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	MGet(ctx context.Context, keys ...string) *redis.SliceCmd
	Ping(ctx context.Context) *redis.StatusCmd
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	HGet(ctx context.Context, key, field string) *redis.StringCmd
	HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd
	HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd
	Expire(ctx context.Context, key string, ttl time.Duration) *redis.BoolCmd
	Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd
}

// redisAdapter реализует internalClient через встраивание *redis.Client.
type redisAdapter struct {
	*redis.Client
}

type client struct {
	rdb           internalClient
	logger        Logger
	timeout       time.Duration
	tracer        trace.Tracer
	scanBatchSize int // размер порции для SCAN в DelByPrefix; задаётся при инициализации
}

// Logger — минимальный интерфейс для логов (реализует github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger).
type Logger interface {
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Warn(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

// NewClient создаёт обёртку над *redis.Client с трейсингом.
// scanBatchSize — размер порции для SCAN в DelByPrefix; должен быть > 0 (при инициализации задаётся снаружи).
func NewClient(rdb *redis.Client, logger Logger, timeout time.Duration, tracerName string, scanBatchSize int) RedisClient {
	return newClient(&redisAdapter{rdb}, logger, timeout, tracerName, scanBatchSize)
}

// newClient создаёт client с внутренним интерфейсом internalClient. scanBatchSize <= 0 заменяется на 100.
func newClient(rdb internalClient, logger Logger, timeout time.Duration, tracerName string, scanBatchSize int) RedisClient {
	if tracerName == "" {
		tracerName = "redis.client"
	}
	if scanBatchSize <= 0 {
		scanBatchSize = 100
	}
	return &client{
		rdb:           rdb,
		logger:        logger,
		timeout:       timeout,
		tracer:        otel.GetTracerProvider().Tracer(tracerName),
		scanBatchSize: scanBatchSize,
	}
}

// withTrace выполняет операцию Redis с трейсингом
func (c *client) withTrace(
	ctx context.Context,
	operation string,
	key string,
	fn func(context.Context) error,
) error {
	ctx, span := c.tracer.Start(ctx, fmt.Sprintf("redis.%s", operation),
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("redis.operation", operation),
			attribute.String("redis.key", key),
		),
	)
	defer span.End()

	traceID := tracing.TraceIDFromContext(ctx)
	spanID := span.SpanContext().SpanID().String()

	start := time.Now()
	err := fn(ctx)
	duration := time.Since(start)

	if err != nil {
		c.logger.Error(ctx, "Redis operation failed",
			zap.String("trace_id", traceID),
			zap.String("span_id", spanID),
			zap.String("operation", operation),
			zap.String("key", key),
			zap.Duration("duration", duration),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(attribute.String("redis.status", "error"))
		return err
	}

	span.SetStatus(codes.Ok, "success")
	span.SetAttributes(
		attribute.String("redis.status", "success"),
		attribute.Float64("redis.duration_us", float64(duration.Microseconds())),
	)

	return nil
}

// withTraceGet выполняет Get операцию с трейсингом и отслеживанием hit/miss
func (c *client) withTraceGet(
	ctx context.Context,
	operation string,
	key string,
	fn func(context.Context) ([]byte, error),
) ([]byte, error) {
	ctx, span := c.tracer.Start(ctx, fmt.Sprintf("redis.%s", operation),
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("redis.operation", operation),
			attribute.String("redis.key", key),
		),
	)
	defer span.End()

	traceID := tracing.TraceIDFromContext(ctx)
	spanID := span.SpanContext().SpanID().String()

	start := time.Now()
	data, err := fn(ctx)
	duration := time.Since(start)

	if err != nil {
		c.logger.Error(ctx, "Redis get failed",
			zap.String("trace_id", traceID),
			zap.String("span_id", spanID),
			zap.String("operation", operation),
			zap.String("key", key),
			zap.Duration("duration", duration),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(attribute.String("redis.status", "error"))
		return nil, err
	}

	// Определяем hit или miss
	isHit := data != nil
	status := "miss"
	if isHit {
		status = "hit"
	}

	span.SetStatus(codes.Ok, status)
	span.SetAttributes(
		attribute.String("redis.status", status),
		attribute.Bool("redis.hit", isHit),
		attribute.Float64("redis.duration_us", float64(duration.Microseconds())),
	)

	return data, nil
}

func (c *client) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.withTrace(ctx, "set", key, func(ctx context.Context) error {
		return c.rdb.Set(ctx, key, value, ttl).Err()
	})
}

func (c *client) Get(ctx context.Context, key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.withTraceGet(ctx, "get", key, func(ctx context.Context) ([]byte, error) {
		data, err := c.rdb.Get(ctx, key).Bytes()
		if errors.Is(err, redis.Nil) {
			return nil, nil // Cache miss, но не ошибка
		}
		return data, err
	})
}

// MGet возвращает значения по ключам одним запросом. Порядок результатов соответствует порядку keys.
func (c *client) MGet(ctx context.Context, keys ...string) ([][]byte, error) {
	if len(keys) == 0 {
		return nil, nil
	}
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	ctx, span := c.tracer.Start(ctx, "redis.mget",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("redis.operation", "mget"),
			attribute.Int("redis.keys_count", len(keys)),
		),
	)
	defer span.End()

	start := time.Now()
	result, err := c.rdb.MGet(ctx, keys...).Result()
	duration := time.Since(start)

	if err != nil {
		c.logger.Error(ctx, "Redis mget failed",
			zap.String("trace_id", tracing.TraceIDFromContext(ctx)),
			zap.Int("keys_count", len(keys)),
			zap.Duration("duration", duration),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	out := make([][]byte, len(result))
	for i, v := range result {
		if v == nil {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		out[i] = []byte(s)
	}

	span.SetStatus(codes.Ok, "success")
	span.SetAttributes(
		attribute.String("redis.status", "success"),
		attribute.Int("redis.keys_count", len(keys)),
		attribute.Float64("redis.duration_us", float64(duration.Microseconds())),
	)
	return out, nil
}

func (c *client) Del(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.withTrace(ctx, "del", key, func(ctx context.Context) error {
		return c.rdb.Del(ctx, key).Err()
	})
}

func (c *client) DelByPrefix(ctx context.Context, prefix string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	pattern := prefix + "*"
	ctx, span := c.tracer.Start(ctx, "redis.delbyprefix",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("redis.operation", "delbyprefix"),
			attribute.String("redis.prefix", prefix),
		),
	)
	defer span.End()

	var cursor uint64
	var deleted int
	for {
		keys, nextCursor, err := c.rdb.Scan(ctx, cursor, pattern, int64(c.scanBatchSize)).Result()
		if err != nil {
			c.logger.Error(ctx, "Redis DelByPrefix scan failed",
				zap.String("prefix", prefix),
				zap.Error(err),
			)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return fmt.Errorf("scan: %w", err)
		}
		if len(keys) > 0 {
			if err := c.rdb.Del(ctx, keys...).Err(); err != nil {
				c.logger.Error(ctx, "Redis DelByPrefix del failed",
					zap.String("prefix", prefix),
					zap.Int("keys_count", len(keys)),
					zap.Error(err),
				)
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				return fmt.Errorf("del: %w", err)
			}
			deleted += len(keys)
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	span.SetStatus(codes.Ok, "success")
	span.SetAttributes(
		attribute.String("redis.status", "success"),
		attribute.Int("redis.deleted_count", deleted),
	)
	return nil
}

func (c *client) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.withTrace(ctx, "ping", "", func(ctx context.Context) error {
		return c.rdb.Ping(ctx).Err()
	})
}

func (c *client) HSet(ctx context.Context, key string, values map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	args := make([]interface{}, 0, len(values)*2)
	for k, v := range values {
		args = append(args, k, v)
	}
	return c.withTrace(ctx, "hset", key, func(ctx context.Context) error {
		return c.rdb.HSet(ctx, key, args...).Err()
	})
}

func (c *client) HGet(ctx context.Context, key, field string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	ctx, span := c.tracer.Start(ctx, "redis.hget",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("redis.operation", "hget"),
			attribute.String("redis.key", key),
			attribute.String("redis.field", field),
		),
	)
	defer span.End()

	traceID := tracing.TraceIDFromContext(ctx)
	spanID := span.SpanContext().SpanID().String()

	start := time.Now()
	result, err := c.rdb.HGet(ctx, key, field).Result()
	duration := time.Since(start)

	if errors.Is(err, redis.Nil) {
		span.SetStatus(codes.Ok, "miss")
		span.SetAttributes(
			attribute.String("redis.status", "miss"),
			attribute.Bool("redis.hit", false),
			attribute.Float64("redis.duration_us", float64(duration.Microseconds())),
		)
		return "", nil
	}

	if err != nil {
		c.logger.Error(ctx, "Redis hget failed",
			zap.String("trace_id", traceID),
			zap.String("span_id", spanID),
			zap.String("key", key),
			zap.String("field", field),
			zap.Duration("duration", duration),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(attribute.String("redis.status", "error"))
		return "", err
	}

	span.SetStatus(codes.Ok, "hit")
	span.SetAttributes(
		attribute.String("redis.status", "hit"),
		attribute.Bool("redis.hit", true),
		attribute.Float64("redis.duration_us", float64(duration.Microseconds())),
	)

	return result, nil
}

func (c *client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	ctx, span := c.tracer.Start(ctx, "redis.hgetall",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("redis.operation", "hgetall"),
			attribute.String("redis.key", key),
		),
	)
	defer span.End()

	traceID := tracing.TraceIDFromContext(ctx)
	spanID := span.SpanContext().SpanID().String()

	start := time.Now()
	result, err := c.rdb.HGetAll(ctx, key).Result()
	duration := time.Since(start)

	if err != nil {
		c.logger.Error(ctx, "Redis hgetall failed",
			zap.String("trace_id", traceID),
			zap.String("span_id", spanID),
			zap.String("key", key),
			zap.Duration("duration", duration),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(attribute.String("redis.status", "error"))
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	span.SetAttributes(
		attribute.String("redis.status", "success"),
		attribute.Int("redis.items_count", len(result)),
		attribute.Float64("redis.duration_us", float64(duration.Microseconds())),
	)

	return result, nil
}

func (c *client) HDel(ctx context.Context, key string, fields ...string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.withTrace(ctx, "hdel", key, func(ctx context.Context) error {
		return c.rdb.HDel(ctx, key, fields...).Err()
	})
}

func (c *client) Expire(ctx context.Context, key string, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.withTrace(ctx, "expire", key, func(ctx context.Context) error {
		return c.rdb.Expire(ctx, key, ttl).Err()
	})
}
