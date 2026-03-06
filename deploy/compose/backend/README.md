# Backend (образ из registry)

Образ собирается в CI (`.github/workflows/build-push.yml`) и пушится в GitHub Container Registry. В образ вшит конфиг `backend/config/production.yaml`; переменные (MYSQL_*, REDIS_*, SESSION_*) подставляются из `env_file`.

## Переменные окружения (на сервере)

В `../../env/<MODE>/.env` или в `.env` рядом с `docker-compose.yml` задайте:

- **BACKEND_IMAGE** — полное имя образа без тега, например `ghcr.io/OWNER/taskflow-backend`
- **IMAGE_TAG** — тег образа: `latest` или короткий SHA коммита (например `a1b2c3d`)
- **MODE** — `local` / `dev` / `prod` (каталог в `deploy/env/`)
- Для конфига приложения: **MYSQL_HOST**, **MYSQL_PORT**, **MYSQL_USER**, **MYSQL_PASSWORD**, **MYSQL_DATABASE**, **REDIS_ADDR**, **REDIS_PASSWORD**, **SESSION_COOKIE_DOMAIN** и др. (см. `backend/config/production.yaml`)

Чтобы подменить конфиг с хоста, раскомментируйте volume в `docker-compose.yml` и положите свой `config/production.yaml`.

## Запуск одного сервиса

```bash
export BACKEND_IMAGE=ghcr.io/ВАШ_ORG/taskflow-backend
export IMAGE_TAG=latest   # или abc1234 после деплоя конкретного коммита

docker compose --env-file ../../env/prod/.env up -d backend
```

## Откат на предыдущую версию

```bash
export IMAGE_TAG=предыдущий_sha   # из истории Actions или логов деплоя
docker compose pull backend
docker compose up -d backend
```
