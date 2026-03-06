# Docker-образы для деплоя (backend, otel, traefik)

Краткая проверка для **локального** развёртывания.

---

## deploy/docker/backend (приложение)

Образ собирается в **CI** (см. `.github/workflows/build-push.yml`) и публикуется в GitHub Container Registry. На сервере используется только **pull + up** (вариант B), без сборки на проде. Подробные шаги — в [deploy/DEPLOY.md](../DEPLOY.md) и [deploy/compose/backend/README.md](../compose/backend/README.md).

## Общие требования для локального запуска

1. **Сеть:** имя сети задаётся по окружению в `.env`: **`DOCKER_NETWORK_NAME=taskflow-local`** (или `taskflow-dev` / `taskflow-prod`).  
   Создание: `docker network create taskflow-local` (или соответствующее окружению). Все сервисы одного окружения должны быть в одной сети.

2. **Переменные:** в `deploy/env/local/.env` (или dev/prod) — переменные из `.env.example`, в т.ч. **`DOCKER_NETWORK_NAME`** для Traefik и всех compose.

---

## deploy/docker/otel (OpenTelemetry Collector)

- **Для локального режима** используется конфиг по умолчанию: `OTEL_CONFIG=/etc/otel-collector-config.yaml` (задаётся в entrypoint, можно переопределить через env).
- **Зависимости по хостам в конфиге (dev):**
  - `elasticsearch:9200` — логи
  - `prometheus:9090` — метрики (remote write)
  - `jaeger:4317` — трейсы (OTLP)
- **Аутентификация Elasticsearch (локально):** через Basic Auth. Либо задать `ELASTICSEARCH_BASIC_AUTH` (base64 `elastic:password`), либо только `ELASTIC_PASSWORD` — entrypoint сгенерирует `ELASTICSEARCH_BASIC_AUTH`.
- **Порты контейнера:** 13133 (health), 4317 (OTLP gRPC), 4318 (OTLP HTTP). Доступ только внутри Docker-сети, проброс на хост не нужен.

Для production используется `otel-collector-config.secure.yaml` (API Key к Elasticsearch, см. `ELASTICSEARCH_API_KEY_SETUP.md`).

---

## deploy/docker/traefik

- **Для локального режима** при запуске контейнера нужно указать конфиг:  
  `--configFile=/etc/traefik/traefik.yml`  
  (образ по умолчанию не выбирает dev/prod сам).
- **Обязательная переменная:** `DOCKER_NETWORK_NAME` — имя сети по окружению (`taskflow-local` / `taskflow-dev` / `taskflow-prod`). Подставляется в `traefik.yml` в `providers.docker.network` и во все compose.
- **Файл маршрутизации (локально):** `dynamic.yml`. В нём зашиты имена сервисов и порты:
  - Frontend: `school_schedule_frontend:3001`
  - Backend: `school_schedule_backend:4000`
  - Grafana: `grafana:3000`
  - Kibana: `kibana:5601`
  - Jaeger: `jaeger:16686`
- **Доступ (локально):** только HTTP (порт 80). Dashboard: `http://traefik.localhost/dashboard/`, API: `http://traefik.localhost/api/`. Остальные сервисы — по правилам из `dynamic.yml` (localhost, api.localhost, metrics.localhost, logs.localhost, jaeger.localhost).

Для production используется `traefik.prod.yml` + `dynamic.prod.yml` (HTTPS, Let's Encrypt, другие имена контейнеров/домены).

---

## Итог проверки

| Компонент | Для локального деплоя |
|-----------|------------------------|
| **OTEL**  | Конфиг `otel-collector-config.yaml` корректен; зависимости: elasticsearch, prometheus, jaeger в одной сети; задать `ELASTIC_PASSWORD` или `ELASTICSEARCH_BASIC_AUTH`. |
| **Traefik** | Задать `DOCKER_NETWORK_NAME` в .env; запускать с `--configFile=/etc/traefik/traefik.yml`; имена сервисов в `dynamic.yml` должны совпадать с теми, что в вашем docker-compose. |

Единого docker-compose, поднимающего frontend, backend, Traefik, OTEL, Elasticsearch, Prometheus, Grafana, Kibana и Jaeger, в репозитории нет — его нужно собрать отдельно или использовать существующий проект с этими сервисами и подключить образы из `deploy/docker/otel` и `deploy/docker/traefik`.
