# Traefik Production Setup

## 📋 Обзор

Production конфигурация Traefik включает:
- ✅ HTTPS с автоматическим получением сертификатов через Let's Encrypt
- ✅ Автоматический редирект HTTP → HTTPS (301)
- ✅ Защита Dashboard через Basic Auth
- ✅ Публичные порты 80 и 443
- ✅ Реальные домены вместо localhost

## 🚀 Быстрый старт

### 1. Подготовка переменных окружения

В `deploy/env/prod/.env` должны быть заданы:

```env
# Traefik
TRAEFIK_CONTAINER_NAME=traefik
TRAEFIK_HTTP_PORT=80
TRAEFIK_HTTPS_PORT=443
TRAEFIK_DASHBOARD_PORT=8080
TRAEFIK_DASHBOARD_DOMAIN=traefik.classplanner.ru
TRAEFIK_CERT_RESOLVER=le

# Let's Encrypt
ACME_EMAIL=your-email@example.com

# Dashboard защита (Basic Auth)
# Формат: user:password (пароль должен быть захеширован)
# Генерация: echo $(htpasswd -nb user password) | sed -e s/\\$/\\$\\$/g
TRAEFIK_DASHBOARD_BASIC_AUTH_USERS=admin:$$apr1$$...

# IP Whitelist (опционально, можно ограничить доступ к dashboard)
TRAEFIK_DASHBOARD_IP_WHITELIST=0.0.0.0/0

# Docker Network (единая для всех режимов)
# ВАЖНО: сеть должна быть создана как внешняя: `docker network create school_net`
DOCKER_NETWORK_NAME=school_net

# Backend
API_DOMAIN=api.classplanner.ru
BACKEND_CONTAINER_NAME=school_schedule_backend
BACKEND_CONTAINER_PORT=4000

# Frontend
FRONTEND_DOMAIN=d.classplanner.ru
FRONTEND_CONTAINER_NAME=school_schedule_frontend
FRONTEND_CONTAINER_PORT=3001
NEXT_PUBLIC_API_URL=https://api.classplanner.ru

# Email (для ссылок в письмах - должен быть тот домен, который видит пользователь)
EMAIL_FRONTEND_URL=https://d.classplanner.ru
```

### 2. Генерация Basic Auth для Dashboard

```bash
# Установить htpasswd (если нет)
# Ubuntu/Debian: sudo apt-get install apache2-utils
# macOS: brew install httpd

# Генерация хеша пароля
htpasswd -nb admin your-password | sed -e s/\\$/\\$\\$/g

# Результат добавить в TRAEFIK_DASHBOARD_BASIC_AUTH_USERS
```

### 3. Запуск

```bash
# Через Taskfile (рекомендуется)
task compose:prod:up

# Или отдельные сервисы
task compose:up:traefik MODE=prod
task compose:up:backend MODE=prod
task compose:up:frontend MODE=prod
```

## 📁 Структура файлов

```
deploy/
├── compose/
│   ├── backend/
│   │   ├── docker-compose.yml          ← DEV (HTTP, localhost)
│   │   └── docker-compose.prod.yml     ← PROD (HTTP+HTTPS, реальные домены)
│   ├── frontend/
│   │   ├── docker-compose.yml          ← DEV
│   │   └── docker-compose.prod.yml     ← PROD
│   └── traefik/
│       ├── docker-compose.yml           ← DEV (порт 80, без HTTPS)
│       └── docker-compose.prod.yml      ← PROD (порты 80+443, Let's Encrypt)
└── docker/traefik/
    ├── traefik.yml                      ← DEV конфигурация
    ├── traefik.prod.yml                 ← PROD конфигурация (HTTPS, Let's Encrypt)
    ├── dynamic.yml                      ← DEV динамическая конфигурация
    └── dynamic.prod.yml                 ← PROD динамическая конфигурация (Basic Auth)
```

## 🔧 Конфигурация

### Backend Labels (PROD)

```yaml
labels:
  # HTTP роутер (для редиректа на HTTPS)
  traefik.http.routers.backend-http.rule: "Host(`api.classplanner.ru`) || PathPrefix(`/api`)"
  traefik.http.routers.backend-http.entrypoints: "web"
  
  # HTTPS роутер (основной)
  traefik.http.routers.backend-secure.rule: "Host(`api.classplanner.ru`) || PathPrefix(`/api`)"
  traefik.http.routers.backend-secure.entrypoints: "websecure"
  traefik.http.routers.backend-secure.tls: "true"
  traefik.http.routers.backend-secure.tls.certresolver: "le"
```

### Frontend Labels (PROD)

```yaml
labels:
  # HTTP роутер (для редиректа на HTTPS)
  traefik.http.routers.frontend.rule: "Host(`d.classplanner.ru`)"
  traefik.http.routers.frontend.entrypoints: "web"
  
  # HTTPS роутер (основной)
  traefik.http.routers.frontend-secure.rule: "Host(`d.classplanner.ru`)"
  traefik.http.routers.frontend-secure.entrypoints: "websecure"
  traefik.http.routers.frontend-secure.tls: "true"
  traefik.http.routers.frontend-secure.tls.certresolver: "le"
```

## 🔒 Безопасность

### Dashboard

- **Basic Auth** - обязателен для production
- **IP Whitelist** - опционально, можно ограничить доступ
- **HTTPS** - dashboard доступен только по HTTPS

### Порты

- **80** - публичный (HTTP, редирект на HTTPS)
- **443** - публичный (HTTPS)
- **8080** - только localhost (Dashboard)

## 📝 Проверка

### Проверить статус

```bash
docker ps | grep traefik
docker logs traefik
```

### Проверить сертификаты

```bash
# Войти в контейнер
docker exec -it traefik sh

# Проверить сертификаты
ls -la /letsencrypt/
cat /letsencrypt/acme.json
```

### Проверить Dashboard

```bash
# Доступ через localhost (требуется Basic Auth)
curl -u admin:password http://localhost:8080/dashboard/
```

## ⚠️ Важно

1. **Перед первым запуском** убедитесь, что домены указывают на IP сервера
2. **ACME_EMAIL** должен быть реальным email для уведомлений Let's Encrypt
3. **Basic Auth** пароль должен быть захеширован (используйте `htpasswd`)
4. **Порты 80 и 443** должны быть открыты в firewall
5. **Let's Encrypt** имеет лимиты на количество запросов (не более 5 сертификатов в неделю для одного домена)

## 🔄 Отличия от DEV

| Параметр | DEV | PROD |
|----------|-----|------|
| HTTPS | ❌ | ✅ |
| Let's Encrypt | ❌ | ✅ |
| HTTP → HTTPS редирект | ❌ | ✅ |
| Dashboard защита | ❌ (insecure) | ✅ (Basic Auth) |
| Порты | 80 (localhost) | 80, 443 (публичные) |
| Домены | localhost | Реальные домены |

---

**Версия:** 1.0.0  
**Дата:** 2025-12-27

