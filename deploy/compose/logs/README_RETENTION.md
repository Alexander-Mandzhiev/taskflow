# Retention политики для observability данных

## Текущая архитектура

**Логи, метрики и трейсы НЕ хранятся в PostgreSQL БД**, а отправляются в специализированные системы:

1. **Логи** → Elasticsearch (через OTEL Collector)
2. **Метрики** → Prometheus (через OTEL Collector)
3. **Трейсы** → Jaeger/Elasticsearch (через OTEL Collector)

## Рекомендуемые retention политики

### 1. Логи (Elasticsearch)

**Рекомендация: 7 дней**

- **Hot phase** (0-1 день): активные индексы, быстрый доступ
- **Warm phase** (1-7 дней): сжатые индексы, медленный доступ
- **Delete phase** (после 7 дней): автоматическое удаление

**Настройка:**

Автоматически через скрипт:
```bash
cd deploy/compose/logs
chmod +x setup-ilm.sh
ELASTIC_PASSWORD=your_password ./setup-ilm.sh
```

Или вручную:
```bash
# Создать ILM политику
curl -X PUT "http://elasticsearch:9200/_ilm/policy/logs-policy" \
  -H 'Content-Type: application/json' \
  -u elastic:${ELASTIC_PASSWORD} \
  -d @elasticsearch-ilm-policy.json

# Создать index template
curl -X PUT "http://elasticsearch:9200/_index_template/logs-template" \
  -H 'Content-Type: application/json' \
  -u elastic:${ELASTIC_PASSWORD} \
  -d @elasticsearch-index-template.json
```

### 2. Метрики (Prometheus)

**Рекомендация: 30-90 дней (в проде можно до года)**

- **Dev**: 30 дней
- **Production**: 90 дней (можно увеличить до 1 года для долгосрочной аналитики)

**Настройка:**
- Уже добавлено в `docker-compose.yaml`: `--storage.tsdb.retention.time=30d`
- В проде: `--storage.tsdb.retention.time=90d`

### 3. Трейсы (Jaeger/Elasticsearch)

**Рекомендация: 3-7 дней**

- **Dev**: 3 дня (in-memory)
- **Production**: 7 дней (в Elasticsearch)

**Настройка:**
- Для Jaeger in-memory: `MEMORY_MAX_TRACES: "10000"` (ограничение по количеству)
- Для Elasticsearch: использовать ILM политику с retention 7 дней

## Почему не в PostgreSQL БД?

1. **Логи** - слишком большой объем, нужна специализированная система поиска
2. **Метрики** - временные ряды, Prometheus оптимизирован для этого
3. **Трейсы** - графы зависимостей, нужна специализированная система

PostgreSQL БД используется только для бизнес-данных (пользователи, расписания, подписки и т.д.)

## Мониторинг использования

- **Elasticsearch**: `GET /_cat/indices/logs-*?v&h=index,store.size,creation.date.string`
- **Prometheus**: `GET /api/v1/status/tsdb` (показывает размер хранилища)
- **Jaeger**: через UI или API
