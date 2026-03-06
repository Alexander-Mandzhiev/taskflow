# Environment Variables

Все переменные окружения для разных режимов запуска: **local / dev / prod**.

## 📁 Структура

```
env/
├── local/
│   ├── .env.example    # Пример для запуска на ноутбуке
│   └── .env            # Реальные значения (в .gitignore)
├── dev/
│   ├── .env.example    # Пример для docker/dev режима (локально)
│   └── .env            # Реальные значения (в .gitignore)
└── prod/
    ├── .env.example    # Пример для облачного сервера
    └── .env            # Реальные значения (в .gitignore)
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
```

### 2. Использование в docker-compose:

```bash
# Через Taskfile (рекомендуется)
task compose:up MODE=local
task compose:up MODE=dev
task compose:up MODE=prod

# Или напрямую (пример: backend)
cd deploy/compose/backend
set ENV_FILE=../../env/local/.env
set DOCKER_NETWORK_NAME=school_local
docker compose --env-file %ENV_FILE% up -d --build
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

### Application Variables

Переменные для приложения внутри контейнера:
- `POSTGRES_HOST` - имя сервиса в docker-compose (для подключения)
- `POSTGRES_PORT` - порт внутри контейнера
- `REDIS_HOST` - имя сервиса в docker-compose
- И т.д.

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

