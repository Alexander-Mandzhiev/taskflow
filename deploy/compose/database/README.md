# MySQL для разработки

Конфигурация вынесена к compose: образ, порты, `mysql.cnf` (charset, InnoDB) — всё в этой папке. Секреты и переменные — в `deploy/env/local/` (см. `.env.example`).

Запуск из корня проекта:

```bash
cd deploy/compose/database && docker compose --env-file ../../env/local/.env up -d
```

Миграции (бэкенд в монорепо): из каталога `backend/` — `goose -dir db/migration mysql "mkk:mkk_secret@tcp(localhost:3306)/mkk?parseTime=true" up`, либо из корня — `goose -dir backend/db/migration mysql "..." up`. Подставьте свои логин/пароль при отличии от `.env.example`.
