# Миграции БД (Goose + MySQL)

## Установка goose

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Убедитесь, что `$GOPATH/bin` или `$HOME/go/bin` в PATH.

## Переменные окружения

- `GOOSE_DRIVER=mysql`
- `GOOSE_DBSTRING="user:password@tcp(host:3306)/dbname?parseTime=true"`  
  Или без env: передавать DSN в каждой команде.

## Команды

Миграции лежат в `backend/db/migration`. Запуск goose:

- **Из каталога `backend/`** (корень бэкенда):

```bash
cd backend
goose -dir db/migration mysql "user:password@tcp(localhost:3306)/mkk?parseTime=true" status
goose -dir db/migration mysql "user:password@tcp(localhost:3306)/mkk?parseTime=true" up
goose -dir db/migration mysql "user:password@tcp(localhost:3306)/mkk?parseTime=true" down
goose -dir db/migration create name_of_migration sql
```

- **Из корня монорепо:**

```bash
goose -dir backend/db/migration mysql "user:password@tcp(localhost:3306)/mkk?parseTime=true" up
```

Пример с env (Linux/macOS), из `backend/`:

```bash
export GOOSE_DRIVER=mysql
export GOOSE_DBSTRING="mkk:mkk_secret@tcp(localhost:3306)/mkk?parseTime=true"
goose -dir db/migration status
goose -dir db/migration up
```

**Локально (Docker):** конфигурация в монорепо — `deploy/compose/database/`, переменные в `deploy/env/local/` (скопировать `.env.example` → `.env`). Поднять MySQL из корня монорепо:

```bash
cd deploy/compose/database && docker compose --env-file ../../env/local/.env up -d
```

Запуск миграций (DSN под `.env.example`), из `backend/`:

```bash
goose -dir db/migration mysql "mkk:mkk_secret@tcp(localhost:3306)/mkk?parseTime=true" up
```

## Автозапуск при старте backend

По умолчанию backend применяет миграции при старте приложения (до поднятия HTTP сервера).
Управляется переменными окружения:

- `MIGRATIONS_AUTO=true|false` — включить/выключить автозапуск (по умолчанию `true`)
- `MIGRATIONS_DIR=/path/to/migration` — путь к каталогу миграций (если не задан, пробуем `./db/migration` и `./backend/db/migration`)

**Без Docker:** создайте базу вручную: `CREATE DATABASE mkk CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;`
