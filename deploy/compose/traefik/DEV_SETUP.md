# Настройка Traefik для DEV режима (без HTTPS)

## 🎯 Цель

Настроить Traefik для dev разработки **без HTTPS редиректа** - только HTTP для простоты.

## 📋 Конфигурация для dev/.env

### Базовые настройки Traefik

```bash
# deploy/env/dev/.env

# ============================================================================
# TRAEFIK - Reverse Proxy (только для dev)
# ============================================================================

# Имя контейнера
TRAEFIK_CONTAINER_NAME=traefik_dev

# Порты
TRAEFIK_HTTP_PORT=80
TRAEFIK_HTTPS_PORT=443
TRAEFIK_DASHBOARD_PORT=8080

# Docker сеть (единая для всех режимов)
# ВАЖНО: сеть должна быть создана как внешняя: `docker network create school_net`
DOCKER_NETWORK_NAME=school_net

# Dashboard
TRAEFIK_DASHBOARD=true
TRAEFIK_DASHBOARD_INSECURE=true  # Прямой доступ через порт 8080 без Basic Auth
TRAEFIK_DASHBOARD_ENABLE=true    # Роутинг через Traefik
TRAEFIK_DASHBOARD_DOMAIN=traefik.localhost  # Или ваш домен
TRAEFIK_DASHBOARD_MIDDLEWARES=  # Пусто для dev (без Basic Auth)
TRAEFIK_DASHBOARD_BASIC_AUTH_USERS=  # Пусто для dev
TRAEFIK_DASHBOARD_IP_WHITELIST=0.0.0.0/0

# HTTP → HTTPS редирект (ОТКЛЮЧЕН для dev)
TRAEFIK_HTTP_REDIRECT=false

# Let's Encrypt (не используется в dev без HTTPS)
ACME_EMAIL=
TRAEFIK_LETSENCRYPT_VOLUME_NAME=traefik_dev_letsencrypt

# ============================================================================
# BACKEND - API
# ============================================================================

BACKEND_CONTAINER_NAME=school_schedule_backend
BACKEND_CONTAINER_PORT=4000
BACKEND_HTTP_PRIORITY=10
BACKEND_HTTP_MIDDLEWARES=backend-stripprefix  # Для PathPrefix роутинга
BACKEND_TLS=false  # Без HTTPS
BACKEND_CERT_RESOLVER=
BACKEND_HTTPS_PRIORITY=10

# Домен для API (если пустой - используется PathPrefix `/api`)
API_DOMAIN=api.localhost  # Или ваш домен, например api.dev.classplanner.ru

# ============================================================================
# FRONTEND - Next.js
# ============================================================================

FRONTEND_CONTAINER_NAME=school_schedule_frontend
FRONTEND_CONTAINER_PORT=3001
FRONTEND_HTTP_PRIORITY=1
FRONTEND_HTTP_MIDDLEWARES=
FRONTEND_TLS=false  # Без HTTPS
FRONTEND_CERT_RESOLVER=
FRONTEND_HTTPS_PRIORITY=10

# Домен для Frontend (если пустой - используется PathPrefix `/`)
FRONTEND_DOMAIN=localhost  # Или ваш домен, например dev.classplanner.ru

# NEXT_PUBLIC_API_URL должен указывать на backend через Traefik
# Если используете домены:
NEXT_PUBLIC_API_URL=http://api.localhost
# Или если используете PathPrefix:
NEXT_PUBLIC_API_URL=http://localhost/api
```

## 🚀 Варианты роутинга

### Вариант 1: Path-based (localhost, без доменов)

**Для локальной разработки без настройки DNS:**

```bash
# .env
API_DOMAIN=  # Пусто - будет использоваться PathPrefix
FRONTEND_DOMAIN=  # Пусто - будет использоваться PathPrefix
NEXT_PUBLIC_API_URL=http://localhost/api
```

**Доступ:**
- Frontend: `http://localhost/`
- Backend API: `http://localhost/api/health`

**Плюсы:**
- ✅ Не требует настройки DNS
- ✅ Работает сразу после запуска
- ✅ Просто для разработки

**Минусы:**
- ❌ Нет HTTPS
- ❌ Нужно настраивать CORS для работы с `/api`

### Вариант 2: Domain-based (с доменами, но без HTTPS)

**Для dev сервера с настроенным DNS:**

```bash
# .env
API_DOMAIN=api.dev.classplanner.ru
FRONTEND_DOMAIN=dev.classplanner.ru
NEXT_PUBLIC_API_URL=http://api.dev.classplanner.ru
TRAEFIK_DASHBOARD_DOMAIN=traefik.dev.classplanner.ru
```

**Доступ:**
- Frontend: `http://dev.classplanner.ru/`
- Backend API: `http://api.dev.classplanner.ru/health`
- Traefik Dashboard: `http://traefik.dev.classplanner.ru/dashboard/` или `http://localhost:8080/dashboard/`

**Плюсы:**
- ✅ Ближе к production окружению
- ✅ Не нужно настраивать CORS
- ✅ Каждый сервис на своем домене

**Минусы:**
- ❌ Требует настройки DNS
- ❌ Нет HTTPS (но можно добавить позже)

## 🔧 Настройка /etc/hosts (для localhost доменов)

Если используете `*.localhost` домены, добавьте в `/etc/hosts`:

```bash
# Windows: C:\Windows\System32\drivers\etc\hosts
# Linux/Mac: /etc/hosts

127.0.0.1 localhost
127.0.0.1 api.localhost
127.0.0.1 traefik.localhost
```

## 📝 Примеры использования

### Запуск dev окружения

```bash
# Через Taskfile
task compose:dev:up

# Или вручную
cd deploy/compose/traefik
docker compose --env-file ../../env/dev/.env up -d

cd ../database
docker compose --env-file ../../env/dev/.env up -d

cd ../cache
docker compose --env-file ../../env/dev/.env up -d

cd ../backend
docker compose --env-file ../../env/dev/.env up -d --build

cd ../frontend
docker compose --env-file ../../env/dev/.env up -d --build
```

### Проверка работы

```bash
# Проверить статус
docker ps | grep traefik
docker ps | grep backend
docker ps | grep frontend

# Проверить роутинг
curl http://localhost/api/health  # Backend через PathPrefix
curl http://api.localhost/health  # Backend через Host (если настроен домен)
curl http://localhost/  # Frontend

# Traefik Dashboard
curl http://localhost:8080/dashboard/  # Прямой доступ
curl http://traefik.localhost/dashboard/  # Через Traefik (если настроен домен)
```

## 🔄 Переключение на HTTPS (если понадобится)

Если позже захотите включить HTTPS:

```bash
# В .env
TRAEFIK_HTTP_REDIRECT=true
API_DOMAIN=api.dev.classplanner.ru
FRONTEND_DOMAIN=dev.classplanner.ru
ACME_EMAIL=admin@classplanner.ru
BACKEND_TLS=true
FRONTEND_TLS=true
TRAEFIK_DASHBOARD_TLS=true
```

## ⚠️ Важно

1. **Production настройки** остаются в `deploy/compose/production/docker-compose.yml` - там нет Traefik
2. **Dev настройки** в отдельных файлах `deploy/compose/*/docker-compose.yml` - они используют Traefik
3. **Переменные** задаются в `deploy/env/dev/.env` для dev режима
4. **NEXT_PUBLIC_API_URL** должен быть задан на этапе build (через `build.args`)

---

**Версия:** 1.0.0  
**Дата:** 2025-01-27

