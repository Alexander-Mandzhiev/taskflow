# Taskflow

**REST API для управления задачами в командах** — сервис с ролевой моделью, историей изменений и поддержкой совместной работы.

## Описание

Taskflow — бэкенд-сервис, который позволяет командам создавать и вести задачи: назначать исполнителей, комментировать, отслеживать изменения и фильтровать по статусу, команде и исполнителю. Реализованы роли в командах (owner, admin, member), аудит изменений задач и кеширование для ускорения типовых запросов.

## Основные возможности

- **Учёт и доступ** — регистрация, вход по JWT, привязка пользователей к командам с ролями
- **Команды** — создание команд, приглашение участников (доступно owner/admin)
- **Задачи** — создание, обновление, фильтрация по `team_id`, `status`, `assignee_id` с пагинацией
- **История** — полная история изменений задачи (кто и когда что изменил)
- **Комментарии** — комментарии к задачам от участников команды
- **Аналитика** — отчёты по командам (участники, закрытые задачи за период), топ авторов задач, проверка целостности (assignee в команде)

## Стек

| Компонент | Технология |
|-----------|------------|
| Язык      | Go         |
| БД        | MySQL      |
| Кеш       | Redis      |
| Инфра     | Docker, Docker Compose |
| Конфиг    | YAML / ENV |

## Запуск

Подробности — в [deploy/docker/README.md](deploy/docker/README.md). Кратко: поднять окружение через Docker Compose, задать переменные по примеру из `deploy/env/local/.env.example`.

## Задачи (Taskfile)

Все команды выполняются из **корня репозитория** через [Task](https://taskfile.dev/): `task <имя>`. Переменные (например `MODE=dev`) задаются при вызове.

### Установка инструментов

| Задача | Описание |
|--------|----------|
| `task install` | Установить все инструменты в `./bin`: gofumpt, gci, golangci-lint, mockery |
| `task install-formatters` | Только форматтеры (gofumpt, gci) |
| `task install-golangci-lint` | Только golangci-lint |
| `task install-mockery` | Только mockery |

### Код: форматирование и линтинг

| Задача | Описание |
|--------|----------|
| `task format` | Форматирование backend: gci + gofumpt |
| `task lint` | Запуск golangci-lint по backend |
| `task check` | format + lint (перед коммитом) |

### Моки

| Задача | Описание |
|--------|----------|
| `task gen` | Сгенерировать моки интерфейсов (mockery по `.mockery.yaml`) |

### Тесты

| Задача | Описание |
|--------|----------|
| `task test` | Все unit-тесты backend (`go test ./...`) |
| `task test:integration` | E2E с testcontainers (MySQL, Redis, backend в Docker). **Нужен запущенный Docker.** Таймаут 10 мин. При ошибке логов в Docker Desktop: в Settings → Docker Engine добавить `"log-driver": "json-file"` и перезапустить. |
| `task test-coverage` | Unit-тесты с покрытием по API, service, adapter; отчёты в `coverage/` |
| `task coverage:func` | Показать покрытие по функциям (после test-coverage) |
| `task coverage:html` | Сгенерировать `coverage/coverage.html` |

### API и smoke

| Задача | Описание |
|--------|----------|
| `task api:smoke` | Smoke-проверка API (бэкенд должен быть поднят). По умолчанию `API_BASE_URL=http://localhost:4000`. Переопределение: `task api:smoke API_BASE_URL=http://localhost:8080` |
| `task api:smoke:build` | Собрать бинарник e2e apiclient в `./bin/apiclient` |

### Docker Compose (окружения)

Режимы: `MODE=local` (по умолчанию), `MODE=dev`, `MODE=prod`. Переменные берутся из `deploy/env/<MODE>/.env`.

| Задача | Описание |
|--------|----------|
| `task compose:local:up` | Поднять инфру для локальной разработки: Traefik, MySQL, Redis, observability. Backend запускается вручную на хосте. |
| `task compose:local:down` | Остановить инфру local |
| `task compose:dev:up` | Всё в Docker, включая backend (сборка образа) |
| `task compose:dev:down` | Остановить dev-окружение |
| `task compose:prod:up` / `compose:prod:down` | Prod: образ backend из registry |
| `task db:up` | Только сеть + MySQL (`MODE` по умолчанию local) |
| `task db:down` | Остановить MySQL |
| `task compose:up:database` | Поднять MySQL (можно указать `MODE=dev`) |
| `task compose:up:cache` | Поднять Redis |
| `task compose:up:backend` | Поднять backend-сервис |
| `task compose:build:backend` | Собрать образ backend локально (тег `taskflow-backend:local`) |

Полный список: `task --list-all`.

## Покрытие тестами

При расчёте покрытия учитывается только **бизнес-логика**: API account v1 и сервис account (пакеты `cmd`, моки, `app` и остальные `pkg` в тотал не входят). Отчёты в `backend/coverage/` (по образцу монорепо; при добавлении модулей — расширить `MODULES` и при необходимости склейку в `total.out`). Команды: `task test-coverage`, `task coverage:html` (файл `backend/coverage/coverage.html`).

## Требования к окружению

- Go 1.25.7
- MySQL 8+
- Redis 6+
- Docker и Docker Compose (для развёртывания)

---

*REST API для управления задачами в командах: роли, история изменений, комментарии, кеширование и аналитика.*
