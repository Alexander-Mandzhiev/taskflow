# MySQL для разработки

Конфигурация вынесена к compose: образ, порты, `mysql.cnf` (charset, InnoDB) — всё в этой папке. Секреты и переменные — в `deploy/env/local/` (см. `.env.example`).

Запуск (сеть задаётся в .env как `DOCKER_NETWORK_NAME=taskflow-local` и т.д.):

```bash
# Сначала создать сеть: docker network create taskflow-local
cd deploy/compose/database && MODE=local docker compose --env-file ../../env/local/.env up -d
```

Миграции (бэкенд в монорепо): из каталога `backend/` — `goose -dir db/migration mysql "mkk:mkk_secret@tcp(localhost:3306)/mkk?parseTime=true" up`, либо из корня — `goose -dir backend/db/migration mysql "..." up`. Подставьте свои логин/пароль при отличии от `.env.example`.

### Предупреждение «World-writable config file ... is ignored»

Если в логах MySQL видно, что `custom.cnf` игнорируется из‑за прав доступа: на Linux/macOS выполните `chmod 644 mysql.cnf` в каталоге `deploy/compose/database/`. На Windows при монтировании файла права внутри контейнера часто выглядят как «world-writable», предупреждение безвредно — MySQL просто использует свои дефолты (подключение и работа БД не страдают).
