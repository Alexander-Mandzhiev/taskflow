# Исправления конфигурации OTEL Collector

## ✅ Исправленные критические проблемы

### 1. API Key аутентификация

**Проблема:** Неправильный синтаксис для API Key

**Исправление:** Используем `headers` для API Key аутентификации:
```yaml
elasticsearch:
  headers:
    Authorization: "ApiKey ${env:ELASTICSEARCH_API_KEY}"
```

**Примечание:** API Key должен быть в формате base64(id:api_key), без префикса "ApiKey".

### 2. CA файл при пустом значении

**Проблема:** Если `ELASTICSEARCH_CA_FILE` не установлен, параметр содержал пустую строку

**Исправление:** 
- Если переменная не установлена, `ca_file` будет пустым (OTEL Collector не будет использовать его)
- **ВАЖНО:** Если `ELASTICSEARCH_CA_FILE` указан, файл **должен существовать**
- Для HTTP (dev) можно не указывать `ELASTICSEARCH_CA_FILE`

### 3. Health check endpoint

**Проблема:** Health check был доступен на `0.0.0.0:13133` (всем)

**Исправление:** Изменен на `localhost:13133` для безопасности
```yaml
extensions:
  health_check:
    endpoint: "localhost:13133"
```

**Примечание:** Docker healthcheck работает через имя сервиса в сети, поэтому это безопасно.

### 4. Traces pipeline

**Проблема:** Traces шли только в `debug`, а не в Elasticsearch

**Исправление:** Добавлен `elasticsearch` в exporters для traces:
```yaml
pipelines:
  traces:
    receivers: [otlp]
    processors: [batch]
    exporters: [elasticsearch]  # Исправлено: было только [debug]
```

### 5. Метрики endpoint

**Статус:** Проверено - метрики доступны через Docker сеть для Prometheus, что безопасно.

## 📝 Рекомендации для production

1. **TLS сертификаты:**
   - Используйте валидные сертификаты (Let's Encrypt или корпоративные)
   - Не используйте `insecure_skip_verify: true` в production

2. **API Key:**
   - Создайте отдельный API Key с минимальными правами (только запись в `logs-*` и `traces-*`)
   - Ротируйте ключи регулярно (например, раз в год)

3. **Мониторинг:**
   - Настройте алерты на ошибки экспорта в Elasticsearch
   - Мониторьте метрики OTEL Collector через Prometheus

4. **Безопасность:**
   - Health check endpoint доступен только локально
   - Метрики доступны только через Docker сеть
   - Используйте TLS для всех подключений к Elasticsearch

## 🔧 Переменные окружения

Обязательные:
```env
ELASTICSEARCH_ENDPOINT=https://elasticsearch:9200
ELASTICSEARCH_API_KEY=<base64_encoded_id:api_key>
```

Опциональные:
```env
ELASTICSEARCH_CA_FILE=/etc/otel/certs/ca.crt  # Только если используете self-signed
ELASTICSEARCH_TLS_INSECURE_SKIP_VERIFY=false   # Только для dev с self-signed
```

## ⚠️ Важные замечания

1. **CA файл:** Если `ELASTICSEARCH_CA_FILE` указан, файл должен существовать по указанному пути
2. **API Key формат:** Должен быть base64(id:api_key), без префикса "ApiKey"
3. **Traces:** Теперь отправляются в Elasticsearch, а не только в debug
4. **Health check:** Доступен только локально для безопасности
