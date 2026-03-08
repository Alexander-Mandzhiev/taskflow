# Environment Variables

Все переменные окружения для разных режимов запуска: **local / dev / prod / test**.

## Стратегия: где хранить секреты

| Окружение | Backend | Compose (MySQL, Redis, observability) |
|-----------|---------|----------------------------------------|
| **Local** | Секреты и настройки в `backend/config/local.yaml`. Запуск на хосте (`go run`) — конфиг по умолчанию, без .env. | `.env` из `deploy/env/local/` — только для контейнеров (database, cache, otel, logs, grafana). |
| **Dev / Prod** | Секреты в `.env`, в YAML только подстановки `${VAR}`. В контейнер передаётся `env_file` при запуске. | Тот же `.env` — compose и backend читают одни переменные. |
| **Test** | Для интеграционных тестов: конфиг приложения `backend/config/test.yaml`; при необходимости переопределение через env (хост/порт к тестовой инфре). | `.env` из `deploy/env/test/` — только для контейнеров MySQL и Redis при запуске `task compose:up:database` / `compose:up:cache` с `MODE=test`. |

Итого: локально пароли можно хранить в YAML бэкенда; в dev/prod — только в .env, без секретов в репозитории; test — отдельная сеть и порты, чтобы не мешать local.

## 📁 Структура

```
env/
├── local/
│   ├── .env.example    # Пример для запуска на ноутбуке
│   └── .env            # Реальные значения (в .gitignore)
├── dev/
│   ├── .env.example    # Пример для docker/dev режима (локально)
│   └── .env            # Реальные значения (в .gitignore)
├── prod/
│   ├── .env.example    # Пример для облачного сервера
│   └── .env            # Реальные значения (в .gitignore)
└── test/
    ├── .env.example    # Пример для интеграционных тестов (db + cache с MODE=test)
    └── .env            # Реальные значения (в .gitignore, опционально — можно копировать из example)
```

## 🚀 Использование

### 1. Создать .env файлы из примеров:

```bash
# local
cd deploy/env/local
cp .env.example .env

# dev
cd ../dev
cp .env.example .env

# prod
cd ../prod
cp .env.example .env

# test (для интеграционных тестов)
cd ../test
cp .env.example .env
```

### 2. Сеть по окружению

Имя Docker-сети задаётся переменной **`DOCKER_NETWORK_NAME`** и должно совпадать во всех compose при одном режиме:

| Окружение | Сеть            |
|-----------|------------------|
| local     | `taskflow-local` |
| dev       | `taskflow-dev`   |
| prod      | `taskflow-prod`  |
| test      | `taskflow-test`  |

Перед первым запуском создайте сеть:  
`docker network create taskflow-local` (или taskflow-dev / taskflow-prod / taskflow-test).

### 3. Режим test (интеграционные тесты)

Для интеграционных тестов поднимают только MySQL и Redis с отдельной сетью и портами:

```bash
# Создать сеть и поднять только db + cache (из корня репозитория)
MODE=test task compose:net
MODE=test task compose:up:database
MODE=test task compose:up:cache
# Запуск интеграционных тестов (backend подключается к localhost:3307 и localhost:6379)
# После тестов:
MODE=test task compose:down:cache
MODE=test task compose:down:database
```

В `backend/config/test.yaml` заданы те же креды (testuser, testpass, testdb); при запуске тестовой инфры через compose подключайтесь к `localhost:3307` (MySQL) и `localhost:6379` (Redis), при необходимости переопределяя хост/порт через env.

### 4. Использование в docker-compose

```bash
# С env-файлом (в .env должна быть DOCKER_NETWORK_NAME=taskflow-local и т.д.)
cd deploy/compose/backend
docker compose --env-file ../../env/local/.env up -d
```

## 📝 Переменные

### Docker Compose Variables

Переменные для оркестрации контейнеров:
- `*_CONTAINER_NAME` - имя контейнера
- `*_HOST_PORT` - порт на хосте
- `*_CONTAINER_PORT` - порт внутри контейнера
- `*_VOLUME_NAME` - имя volume
- `*_MEMORY_LIMIT` - лимит памяти
- `*_CPU_LIMIT` - лимит CPU

### Application Variables (dev/prod — backend в контейнере)

В dev/prod backend получает переменные из `env_file`; в YAML (`production.yaml`) только подстановки `${VAR}`:
- `MYSQL_HOST`, `MYSQL_PORT`, `MYSQL_USER`, `MYSQL_PASSWORD`, `MYSQL_DATABASE` — MySQL
- `REDIS_ADDR`, `REDIS_PASSWORD` — Redis
- `DOCKER_NETWORK_NAME` — имя сети (`taskflow-local` / `taskflow-dev` / `taskflow-prod`)

Для **local** при запуске backend на хосте конфиг по умолчанию — `backend/config/local.yaml` (секреты там).

### Frontend build-time variables (ВАЖНО)

Next.js переменные `NEXT_PUBLIC_*` (например `NEXT_PUBLIC_API_URL`) **вшиваются на этапе сборки**.
Поэтому они должны быть:
- в `.env` (который читает `docker compose --env-file ...` для интерполяции)
- и прокинуты в сборку через `build.args` (это уже сделано в `deploy/compose/frontend/docker-compose.yml`)

#### ⚠️ Важно для dev режима:

В **dev режиме** порты **не пробрасываются** наружу - всё идёт через Traefik:
- ✅ `NEXT_PUBLIC_API_URL=https://api.classplanner.ru` - фронтенд обращается к бэкенду через Traefik
- Нужно добавить в `/etc/hosts`: `127.0.0.1 classplanner.ru api.classplanner.ru traefik.classplanner.ru`
- Traefik проксирует запросы по доменам с TLS (Let's Encrypt)
- Backend и Frontend доступны только через Traefik на портах 80/443

**Для local режима** (ноутбук, backend/frontend НЕ в Docker):
- ✅ `NEXT_PUBLIC_API_URL=http://localhost:4000` - фронтенд и бэкенд запускаются на ноутбуке напрямую
- В Docker контейнерах только: db, cache, observability
- Backend и Frontend запускаются вручную на ноутбуке (например, `npm run dev` для frontend, `go run` для backend)

## ⚠️ Важно

- `.env` файлы в `.gitignore` - не коммитятся в git
- `.env.example` файлы в git - примеры для всех
- **Обязательно** измените все `CHANGE_ME_*` значения на реальные
- **Обязательно** используйте сильные пароли для production

---

**Версия:** 1.0.0  
**Дата:** 2025-01-27

