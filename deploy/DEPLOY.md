# Деплой backend по варианту B (образ из CI → registry → сервер)

## Схема

1. **Разработка**: коммит и push в `main`.
2. **CI**: линтинг и тесты (`.github/workflows/ci.yml`).
3. **Build**: при push в `main` собирается Docker-образ и пушится в GitHub Container Registry (`.github/workflows/build-push.yml`). Теги: короткий SHA коммита и `latest`.
4. **Сервер**: подтягиваете образ по тегу и перезапускаете только backend (без сборки на сервере).

---

## Шаг 1. Включить публикацию в GHCR

- В репозитории GitHub: **Settings → Actions → General** — разрешить workflow доступ к **Read and write permissions** для пакетов (или оставить по умолчанию, если уже есть).
- После первого успешного run workflow образ появится в **Packages** репозитория.
- Для **private** репо: на сервере нужна авторизация в GHCR (см. шаг 4).

---

## Шаг 2. Имя образа

После первого push образ будет:

- **Адрес**: `ghcr.io/<OWNER>/taskflow-backend`, например `ghcr.io/Alexander-Mandzhiev/taskflow-backend`
- **Теги**: `latest`, а также короткий SHA коммита (например `a1b2c3d`), под которым прошёл build.

---

## Шаг 3. Подготовка на сервере

1. Установить Docker и Docker Compose.
2. Создать сеть по окружению (если ещё нет):  
   `docker network create taskflow-local` (или `taskflow-dev` / `taskflow-prod`).  
   Имя задаётся в `.env` переменной `DOCKER_NETWORK_NAME`.
3. Склонировать репозиторий (или только каталог `deploy/` и конфиги).
4. В `deploy/env/prod/.env` (или в своём `.env`) задать:
   - переменные приложения (MySQL, Redis, сессии и т.д.);
   - `BACKEND_IMAGE=ghcr.io/<OWNER>/taskflow-backend`;
   - при деплое по SHA: `IMAGE_TAG=<short_sha>` или для последнего билда: `IMAGE_TAG=latest`.
5. При необходимости подменить конфиг: раскомментировать volume в `deploy/compose/backend/docker-compose.yml` и примонтировать свой `production.yaml` (иначе используется конфиг, вшитый в образ).

---

## Шаг 4. Авторизация на сервере в GHCR (для private репо)

```bash
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin
```

`GITHUB_TOKEN` — Personal Access Token с правом `read:packages`. Для публичного репо pull без логина может работать без токена.

---

## Taskfile (локальный подъём и сеть)

Сеть создаётся таской по окружению: `task compose:net` создаёт `taskflow-<MODE>` (local/dev/prod). Запуск сервисов через Task:

```bash
task db:up              # сеть + MySQL (MODE=local по умолчанию)
task compose:up:cache   # Redis
task compose:up:backend # backend из registry
task compose:build:backend  # локальная сборка образа (тег :local)
```

С указанием окружения: `task db:up MODE=dev`, `task compose:up:backend MODE=prod`. В `deploy/env/<MODE>/.env` должен быть задан `DOCKER_NETWORK_NAME=taskflow-<MODE>`.

---

## Шаг 5. Деплой и откат

**Деплой (после push в main):**

```bash
cd deploy/compose/backend
export IMAGE_TAG=latest   # или конкретный sha, например a1b2c3d
docker compose --env-file ../../env/prod/.env pull backend
docker compose --env-file ../../env/prod/.env up -d backend
```

**Только обновить backend (остальные сервисы не трогаются):**

```bash
docker compose pull backend
docker compose up -d backend
```

**Откат на предыдущую версию:**

```bash
export IMAGE_TAG=предыдущий_sha   # из логов деплоя или Actions
docker compose pull backend
docker compose up -d backend
```

---

## Семантические теги (v*)

При push тега вида `v1.2.3` (например `git tag v1.2.3 && git push origin v1.2.3`) workflow Build and Push запускается и пушит образ с тегом **`v1.2.3`** (без `latest`). На сервере можно зафиксировать версию: `IMAGE_TAG=v1.2.3`.
