# Модуль Team

Верхнеуровневое описание репозиторного слоя и моделей.

## Модели (`model/`)

| Модель         | Описание |
|----------------|----------|
| `Team`         | Команда: id, name, created_by, created_at, updated_at, deleted_at |
| `TeamMember`   | Участник: id, user_id, team_id, role, created_at |
| `TeamWithRole`     | Команда + роль текущего пользователя в ней (для списка «мои команды») |
| `TeamWithMembers`  | Команда + список участников (возвращает GetByID) |
| `TeamInput`        | Входные данные: name (создание/обновление команды) |

Роли: `RoleOwner`, `RoleAdmin`, `RoleMember` (константы).

Ошибки: `ErrTeamNotFound`, `ErrMemberNotFound`, `ErrAlreadyMember`, `ErrNilInput`.

## Контракты репозитория

### Низкий уровень (реализуют слой team/ и member/)

- **TeamReaderRepository**: `GetByID(teamID)`, `ListByUserID(userID)` → команды с ролью пользователя (один JOIN).
- **TeamWriterRepository**: `Create(input, createdBy)` → запись в `teams`.
- **MemberReaderRepository**: `GetByTeamID(teamID)`, `GetMember(teamID, userID)`.
- **MemberWriterRepository**: `AddMember(teamID, userID, role)`.

### Верхний уровень (реализует адаптер)

- **TeamRepository** — единая точка доступа для сервиса:
  - `Create(input, ownerUserID)` — создаёт команду и добавляет создателя как owner в `team_members`.
  - `GetByID(teamID)` — команда с участниками (`TeamWithMembers`).
  - `ListByUserID(userID)`, `GetMember(teamID, userID)`, `AddMember(teamID, userID, role)`.

Все методы принимают `(ctx, tx, ...)`. При `tx != nil` — работа в транзакции; при `tx == nil` — вне транзакции.

## Таблицы БД

- `teams` — одна запись на команду.
- `team_members` — связь пользователь–команда с ролью (многие-ко-многим).

Чтение: GetByID внутри делает два запроса (команда + участники). Список «мои команды» — один запрос (teams JOIN team_members по user_id).
