# Локальные переменные окружения

Скопируйте `.env.example` в `.env` и при необходимости отредактируйте:

```bash
cp .env.example .env
```

Файл `.env` подключается в compose (`deploy/compose/database/`): переменные пробрасываются в контейнер и используются при подстановке (например, порт). Не коммитить в git.

Запуск с этим env из корня монорепо (mkk):

```bash
cd deploy/compose/database && docker compose --env-file ../../env/local/.env up -d
```
