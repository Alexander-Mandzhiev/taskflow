# Grafana Dashboards

Предустановленные дашборды для визуализации метрик приложения.

## HTTP Metrics Dashboard

**Файл:** `http-metrics.json`

Дашборд для мониторинга HTTP метрик приложения, собираемых через OpenTelemetry middleware.

### Метрики

1. **HTTP Requests per Second** - количество запросов в секунду (rate)
2. **Current RPS** - текущее значение RPS
3. **Total Requests (1h)** - общее количество запросов за последний час
4. **HTTP Response Time** - график времени ответа (average, p95, p99)
5. **Avg Response Time** - среднее время ответа
6. **p95 Response Time** - 95-й перцентиль времени ответа
7. **HTTP Error Rate** - график ошибок (4xx/5xx) и успешных запросов
8. **Error Rate %** - процент ошибок от общего количества запросов
9. **Errors/sec** - количество ошибок в секунду

### Используемые метрики Prometheus

- `school_schedule_monolith_http_requests_total` - счетчик входящих запросов
- `school_schedule_monolith_http_responses_total` - счетчик ответов (с labels: status, method, path)
- `school_schedule_monolith_http_response_time_seconds` - гистограмма времени ответа (с labels: status)

### Формат метрик

Метрики экспортируются через OpenTelemetry Collector в Prometheus и имеют формат:
- Namespace: `school_schedule`
- App Name: `monolith`
- Полное имя: `{namespace}_{app_name}_{metric_name}`

Гистограммы автоматически преобразуются Prometheus в:
- `{metric_name}_bucket` - buckets для квантилей
- `{metric_name}_count` - количество измерений
- `{metric_name}_sum` - сумма значений

### Обновление дашборда

Дашборд автоматически загружается при старте Grafana через provisioning.
Для применения изменений перезапустите контейнер Grafana:

```bash
docker restart grafana
```

### Добавление новых дашбордов

1. Создайте JSON файл дашборда в этой папке
2. Экспортируйте дашборд из Grafana UI: Dashboard Settings → JSON Model
3. Сохраните как `{dashboard-name}.json`
4. Перезапустите Grafana

