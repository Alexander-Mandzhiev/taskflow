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

	"mkk/pkg/tracing"
)

type client struct {
	rdb     redis.Cmdable
	logger  Logger
	timeout time.Duration
	tracer  trace.Tracer
}

// Logger — минимальный интерфейс для логов (реализует mkk/pkg/logger).
type Logger interface {
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Warn(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

// NewClient создаёт обёртку над go-redis клиентом (ClusterClient или Client) с трейсингом.
// tracerName — название модуля для трейсинга (например "mkk.cache"). Если пустой, используется "redis.client".
func NewClient(rdb redis.Cmdable, logger Logger, timeout time.Duration, tracerName string) RedisClient {
	if tracerName == "" {
		tracerName = "redis.client"
	}
	return &client{
		rdb:     rdb,
		logger:  logger,
		timeout: timeout,
		tracer:  otel.GetTracerProvider().Tracer(tracerName),
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

	c.logger.Debug(ctx, "Redis operation success",
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
		zap.String("operation", operation),
		zap.String("key", key),
		zap.Duration("duration", duration),
	)

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
		c.logger.Debug(ctx, "Redis hit",
			zap.String("trace_id", traceID),
			zap.String("span_id", spanID),
			zap.String("key", key),
			zap.Duration("duration", duration),
		)
	} else {
		c.logger.Debug(ctx, "Redis miss",
			zap.String("trace_id", traceID),
			zap.String("span_id", spanID),
			zap.String("key", key),
			zap.Duration("duration", duration),
		)
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
	return c.withTrace(ctx, "hset", key, func(ctx context.Context) error {
		return c.rdb.HSet(ctx, key, values).Err()
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
		c.logger.Debug(ctx, "Redis hget miss",
			zap.String("trace_id", traceID),
			zap.String("span_id", spanID),
			zap.String("key", key),
			zap.String("field", field),
			zap.Duration("duration", duration),
		)
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

	c.logger.Debug(ctx, "Redis hget hit",
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
		zap.String("key", key),
		zap.String("field", field),
		zap.Duration("duration", duration),
	)

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

	c.logger.Debug(ctx, "Redis hgetall",
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
		zap.String("key", key),
		zap.Int("items", len(result)),
		zap.Duration("duration", duration),
	)

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
