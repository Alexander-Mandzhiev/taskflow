# Настройка API Key для Elasticsearch в OTEL Collector

## Обзор

Использование API Key вместо BasicAuth для подключения OTEL Collector к Elasticsearch - более безопасный подход.

## Преимущества API Key

1. **Гранулярные права** - можно выдать доступ только на запись в индексы логов
2. **Легкая ротация** - можно регенерировать без смены основного пароля
3. **Точечное отключение** - при компрометации можно отключить конкретный ключ
4. **Аудит** - действия привязаны к конкретному API Key

## Шаг 1: Создание API Key в Elasticsearch

### Вариант A: Через Kibana UI

1. Зайдите в Kibana → Stack Management → Security → API Keys
2. Нажмите "Create API key"
3. Укажите:
   - **Name**: `otel-collector-logs`
   - **Expiration**: по необходимости (или без срока)
   - **Role descriptors**: выберите роль с правами на запись в `logs-*` индексы

### Вариант B: Через Elasticsearch API

```bash
# Создание API Key с ограниченными правами
curl -X POST "https://elasticsearch:9200/_security/api_key" \
  -u "elastic:YOUR_PASSWORD" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "otel-collector-logs",
    "role_descriptors": {
      "log_writer": {
        "cluster": ["monitor"],
        "index": [
          {
            "names": ["logs-*"],
            "privileges": ["create_index", "write", "manage"]
          }
        ]
      }
    }
  }'
```

Ответ будет содержать `id` и `api_key`:
```json
{
  "id": "VFR2WU41VUJIbG9SbGJUdVFrMFk",
  "name": "otel-collector-logs",
  "api_key": "NVVhVDE3SDlSQS0wM1Rxb24xdXFldw=="
}
```

**Важно**: Сохраните полный API Key в формате `id:api_key` (base64):
```
VFR2WU41VUJIbG9SbGJUdVFrMFk6NVVhVDE3SDlSQS0wM1Rxb24xdXFldw==
```

## Шаг 2: Настройка переменных окружения

Добавьте в `.env` файл для соответствующего окружения:

```env
# Elasticsearch для OTEL Collector
ELASTICSEARCH_ENDPOINT=https://elasticsearch:9200
# API Key должен быть в формате base64(id:api_key)
# БЕЗ префикса "ApiKey " - он добавляется автоматически в конфиге
ELASTICSEARCH_API_KEY=VFR2WU41VUJIbG9SbGJUdVFrMFk6NVVhVDE3SDlSQS0wM1Rxb24xdXFldw==

# TLS настройки (опционально)
# Для HTTPS с валидным сертификатом - не указывайте ELASTICSEARCH_CA_FILE
# Для self-signed сертификата укажите путь к CA файлу:
ELASTICSEARCH_CA_FILE=/etc/otel/certs/ca.crt
# ВАЖНО: Если ELASTICSEARCH_CA_FILE указан, файл должен существовать!
# Для dev с self-signed можно использовать insecure (не рекомендуется для prod):
# ELASTICSEARCH_TLS_INSECURE_SKIP_VERIFY=true
```

**Формат API Key:**
- Elasticsearch возвращает `id` и `api_key` отдельно
- Нужно закодировать в base64: `echo -n "id:api_key" | base64`
- Результат сохранить в `ELASTICSEARCH_API_KEY` (без префикса "ApiKey")

**Важно для TLS:**
- Если `ELASTICSEARCH_CA_FILE` не указан, `ca_file` будет пустым (OTEL Collector не будет использовать его)
- Если `ELASTICSEARCH_CA_FILE` указан, файл **должен существовать** по указанному пути
- Для HTTP (dev) можно не указывать `ELASTICSEARCH_CA_FILE`

## Шаг 3: Настройка TLS сертификатов (если используется HTTPS)

### Для self-signed сертификатов:

1. Получите CA сертификат от Elasticsearch:
```bash
# Из контейнера Elasticsearch
docker cp elasticsearch:/usr/share/elasticsearch/config/certs/ca/ca.crt ./certs/ca.crt
```

2. Убедитесь, что сертификат доступен в OTEL Collector:
   - Файл должен быть в директории, указанной в `OTEL_CERTS_DIR` (по умолчанию `./certs`)
   - Или укажите полный путь в `ELASTICSEARCH_CA_FILE`

### Для валидных сертификатов:

Если используете валидный сертификат (например, от Let's Encrypt), просто укажите путь к CA bundle:
```env
ELASTICSEARCH_CA_FILE=/etc/ssl/certs/ca-certificates.crt
```

## Шаг 4: Использование secure конфига

Убедитесь, что OTEL Collector использует secure конфиг:

```env
OTEL_CONFIG=/etc/otel-collector-config.secure.yaml
```

Или в docker-compose:
```yaml
environment:
  - OTEL_CONFIG=/etc/otel-collector-config.secure.yaml
```

## Шаг 5: Проверка подключения

Проверьте логи OTEL Collector:
```bash
docker logs otel-collector
```

Должны быть сообщения об успешной отправке логов в Elasticsearch.

Проверьте индексы в Elasticsearch:
```bash
curl -X GET "https://elasticsearch:9200/_cat/indices/logs-*" \
  -u "elastic:YOUR_PASSWORD"
```

## Создание роли для API Key (рекомендуется)

Вместо использования суперпользователя `elastic`, создайте отдельную роль:

```bash
# Создание роли для записи логов
curl -X PUT "https://elasticsearch:9200/_security/role/log_writer" \
  -u "elastic:YOUR_PASSWORD" \
  -H "Content-Type: application/json" \
  -d '{
    "cluster": ["monitor"],
    "indices": [
      {
        "names": ["logs-*"],
        "privileges": ["create_index", "write", "manage"]
      }
    ]
  }'

# Создание пользователя для OTEL Collector
curl -X POST "https://elasticsearch:9200/_security/user/otel_collector" \
  -u "elastic:YOUR_PASSWORD" \
  -H "Content-Type: application/json" \
  -d '{
    "password": "STRONG_PASSWORD_HERE",
    "roles": ["log_writer"]
  }'

# Создание API Key для этого пользователя
curl -X POST "https://elasticsearch:9200/_security/api_key" \
  -u "elastic:YOUR_PASSWORD" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "otel-collector-logs",
    "expiration": "365d"
  }'
```

## Безопасность

1. **Не коммитьте API Key в git** - используйте `.env` файлы, которые в `.gitignore`
2. **Ротируйте ключи регулярно** - например, раз в год
3. **Используйте минимальные права** - только то, что нужно для записи логов
4. **Включите TLS** - даже внутри Docker сети для дополнительной защиты
5. **Мониторьте использование** - проверяйте логи Elasticsearch на подозрительную активность

## Troubleshooting

### Ошибка: "unable to authenticate user"

- Проверьте правильность API Key (должен быть в формате `id:api_key`)
- Убедитесь, что API Key не истек
- Проверьте права API Key

### Ошибка: "certificate verify failed"

- Убедитесь, что `ELASTICSEARCH_CA_FILE` указывает на правильный файл
- Для self-signed сертификатов можно временно использовать `insecure_skip_verify: true` (только для dev)

### Ошибка: "index not found"

- Проверьте права API Key - должно быть право `create_index` для `logs-*`
- Проверьте формат `logs_index` в конфиге OTEL Collector
