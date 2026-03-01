# Redis Client с трейсингом

Redis клиент с автоматическим трейсингом всех операций через OpenTelemetry.

## Особенности

- ✅ Автоматический трейсинг всех операций Redis
- ✅ Метрики hit/miss для Get операций
- ✅ Логирование всех операций с trace_id и span_id
- ✅ Измерение времени выполнения операций
- ✅ Поддержка названия модуля для трейсинга

## Использование

### Создание клиента

**Через BuildClient (рекомендуется)** — создаёт go-redis клиент и обёртку с трейсингом. Настройка через опции (как в pkg/logger, pkg/metric, pkg/tracing):

```go
import "mkk/pkg/cache"

client, err := cache.BuildClient(logger, "mkk.cache",
    cache.WithAddr("localhost:6379"),
    cache.WithDialTimeout(5*time.Second),
)

// С паролем:
client, err = cache.BuildClient(logger, "mkk.cache",
    cache.WithAddr("localhost:6379"),
    cache.WithPassword(redisPassword),
)
```

**Низкоуровнево (если уже есть go-redis клиент):**

```go
import "mkk/pkg/cache"

// rdb — *redis.Client или redis.ClusterClient
client := cache.NewClient(rdb, logger, 5*time.Second, "mkk.cache")
```

### Параметры

- `rdb` — go-redis клиент (Client или ClusterClient)
- `logger` — реализация cache.Logger (например logger.Logger())
- `timeout` — таймаут для операций
- `tracerName` — название модуля для трейсинга (например "mkk.cache")
  - Если пустой, используется "redis.client" по умолчанию
  - Используется для группировки трейсов в Jaeger

## Трейсинг

Все операции автоматически создают spans в OpenTelemetry с атрибутами:

### Для всех операций:
- `redis.operation` = "get" | "set" | "del" | "hget" | "hset" | "hgetall" | "hdel" | "expire" | "ping"
- `redis.key` = ключ Redis
- `redis.status` = "success" | "error" | "hit" | "miss"
- `redis.duration_us` = время выполнения в микросекундах

### Для Get операций дополнительно:
- `redis.hit` = true/false
- `redis.field` = поле для HGet операций
- `redis.items_count` = количество элементов для HGetAll

## Метрики

Метрики доступны через OpenTelemetry и могут быть экспортированы в Prometheus:
- Redis hit/miss ratio
- Redis operation duration
- Redis error rate

## Примеры трейсов в Jaeger

```
redis.get
  - redis.operation: get
  - redis.key: teachers:hash:123
  - redis.status: hit
  - redis.hit: true
  - redis.duration_us: 1234

redis.hget
  - redis.operation: hget
  - redis.key: teachers:hash:123
  - redis.field: 456
  - redis.status: miss
  - redis.hit: false
  - redis.duration_us: 567
```

## Преимущества

1. **Автоматический трейсинг** - не нужно оборачивать операции вручную
2. **Метрики hit/miss** - автоматическое отслеживание эффективности кэша
3. **Единообразие** - все операции Redis трейсятся одинаково
4. **Гибкость** - можно создавать отдельные клиенты для разных модулей

