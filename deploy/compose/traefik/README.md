# Traefik Reverse Proxy

Traefik используется как reverse proxy и load balancer **в local и dev режимах**.

⚠️ **Важно**: 
- В **local режиме** Traefik нужен для доступа к observability через домены (logs.localhost, metrics.localhost, jaeger.localhost)
- В **dev режиме** Traefik используется для всех сервисов (frontend, backend, observability)
- В **production режиме** Traefik не используется. На сервере должен быть свой reverse proxy (nginx, cloudflare, или другой)

## 📁 Структура файлов

```
deploy/
├── docker/traefik/
│   ├── Dockerfile       # Dockerfile для сборки образа
│   ├── traefik.yml     # Статическая конфигурация Traefik (YAML)
│   └── dynamic.yml     # Динамическая конфигурация Traefik (YAML)
└── compose/traefik/
    ├── docker-compose.yml  # Docker Compose конфигурация
    └── README.md           # Документация
```

## 🔧 Конфигурация

### Статическая конфигурация (traefik.yml)

Основная конфигурация находится в файле `deploy/docker/traefik/traefik.yml`:
- **Логирование** - уровень и формат логов
- **API и Dashboard** - настройки веб-интерфейса
- **Entrypoints** - точки входа (порты 80, 443)
- **Providers** - провайдеры конфигурации (Docker)

### Динамическая конфигурация (Docker Labels)

Маршрутизация настраивается через Docker labels в docker-compose файлах сервисов:
- Frontend: `deploy/compose/frontend/docker-compose.yml`
- Backend: `deploy/compose/backend/docker-compose.yml`

## 🚀 Использование

### Запуск через Taskfile

```bash
# Запустить Traefik для local режима (для observability)
task compose:local:up

# Запустить Traefik для dev режима (для всех сервисов)
task compose:dev:up

# Остановить
task compose:local:down  # или compose:dev:down
```

### Переменные окружения

Все переменные задаются в `.env` файле для каждого режима:
- `deploy/env/dev/.env` - для dev режима
- `deploy/env/local/.env` - для local режима
- `deploy/env/prod/.env` - для prod режима

**Основные переменные:**

```bash
# Docker Network
# ВАЖНО: сеть должна быть создана как внешняя: `docker network create school_net`
DOCKER_NETWORK_NAME=school_net

# Traefik
TRAEFIK_CONTAINER_NAME=traefik_dev
TRAEFIK_HTTP_PORT=80
TRAEFIK_DASHBOARD_PORT=8080
TRAEFIK_DASHBOARD=true
TRAEFIK_DASHBOARD_INSECURE=true
TRAEFIK_DASHBOARD_DOMAIN=traefik.localhost
```

## 📋 Настройка роутинга

### Frontend

В `deploy/compose/frontend/docker-compose.yml`:

```yaml
labels:
  - traefik.enable=true
  - traefik.http.services.frontend.loadbalancer.server.port=3001
  - traefik.http.routers.frontend-http.rule=Host(`localhost`) || PathPrefix(`/`)
  - traefik.http.routers.frontend-http.entrypoints=web
  - traefik.http.routers.frontend-http.service=frontend
```

### Backend

В `deploy/compose/backend/docker-compose.yml`:

```yaml
labels:
  - traefik.enable=true
  - traefik.http.services.backend.loadbalancer.server.port=4000
  - traefik.http.routers.backend-http.rule=Host(`api.localhost`) || PathPrefix(`/api`)
  - traefik.http.routers.backend-http.entrypoints=web
  - traefik.http.routers.backend-http.service=backend
  - traefik.http.routers.backend-http.middlewares=backend-stripprefix
  - traefik.http.middlewares.backend-stripprefix.stripprefix.prefixes=/api
```

## 🔍 Отладка

### Проверка конфигурации

```bash
# Проверить, что Traefik видит конфигурацию
docker exec traefik cat /etc/traefik/traefik.yml

# Или посмотреть исходный файл
cat deploy/docker/traefik/traefik.yml

# Проверить логи
docker logs traefik --tail 50

# Проверить API
curl http://localhost:8080/api/rawdata | jq
```

### Traefik Dashboard

Откройте в браузере: `http://localhost:8080/dashboard/`

Здесь можно увидеть:
- **HTTP Routers** - все настроенные роутеры
- **HTTP Services** - все сервисы
- **HTTP Middlewares** - все middleware

## ⚙️ Изменение конфигурации

### Изменить статическую конфигурацию

Отредактируйте `deploy/docker/traefik/traefik.yml` или `deploy/docker/traefik/dynamic.yml` и пересоберите образ:

```bash
# Пересобрать образ и перезапустить
task compose:up:traefik

# Или вручную
cd deploy/compose/traefik
docker compose --env-file ../../env/dev/.env up -d --build
```

### Изменить динамическую конфигурацию (роутинг)

Отредактируйте labels в docker-compose файлах сервисов и перезапустите сервис:

```bash
docker restart school_schedule_frontend_local
```

Traefik автоматически обнаружит изменения (watch: true).

## 🔒 Безопасность для Production

Для production нужно:

1. **Отключить insecure API:**
   ```yaml
   api:
     dashboard: true
     insecure: false  # Требовать авторизацию
   ```

2. **Включить HTTPS:**
   ```yaml
   entryPoints:
     websecure:
       address: ":443"
   certificatesResolvers:
     le:
       acme:
         email: "admin@example.com"
         storage: /letsencrypt/acme.json
         httpChallenge:
           entryPoint: web
   ```

3. **Добавить Basic Auth для Dashboard:**
   ```yaml
   # В labels сервиса traefik
   - traefik.http.middlewares.traefik-auth.basicauth.users=admin:$$apr1$$...
   ```

## 📚 Дополнительные ресурсы

- [Официальная документация Traefik](https://doc.traefik.io/traefik/)
- [Traefik Docker Provider](https://doc.traefik.io/traefik/providers/docker/)
- [Traefik Routing](https://doc.traefik.io/traefik/routing/routers/)
