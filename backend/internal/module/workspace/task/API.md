# API модуля задач (tasks)

Ручки для управления задачами по ТЗ и контракту `TaskService`. Авторизация: JWT в заголовке (userID из токена).

---

## 1. Управление задачами (обязательно по ТЗ)

### 1.1 Создать задачу

- **Метод:** `POST /api/v1/tasks`
- **Сервис:** `TaskService.Create(ctx, userID, teamID, input)`
- **Тело запроса:** `team_id` (UUID) + поля задачи:
  ```json
  {
    "team_id": "uuid",
    "title": "string",
    "description": "string",
    "status": "todo | in_progress | done",
    "assignee_id": "uuid | null"
  }
  ```
- **Ответ:** `201` + объект задачи; при `nil` body — `400` + `ErrNilInput`; недопустимый статус — `400` + `ErrInvalidStatus`; assignee не в команде — `400` + `ErrAssigneeNotInTeam`; пользователь не в команде — `403` + `ErrForbidden`.

---

### 1.2 Список задач с фильтрацией и пагинацией

- **Метод:** `GET /api/v1/tasks?team_id=...&status=...&assignee_id=...&limit=...&offset=...`
- **Сервис:** `TaskService.List(ctx, userID, filter, pagination)`
- **Query:** `team_id`, `status`, `assignee_id` (опционально), `limit`, `offset` (пагинация на уровне БД по ТЗ).
- **Ответ:** `200` + тело вида `{ "items": [...], "total": N }`. Если задан `team_id`, userID должен быть в этой команде; иначе `403`.

---

### 1.3 Получить задачу по ID

- **Метод:** `GET /api/v1/tasks/{id}`
- **Сервис:** `TaskService.GetByID(ctx, taskID, userID)`
- **Ответ:** `200` + объект задачи; не найдена или пользователь не в команде — `404` / `403`.

---

### 1.4 Обновить задачу

- **Метод:** `PUT /api/v1/tasks/{id}`
- **Сервис:** `TaskService.Update(ctx, userID, taskID, input)`
- **Тело запроса:**
  ```json
  {
    "title": "string",
    "description": "string",
    "status": "todo | in_progress | done",
    "assignee_id": "uuid | null"
  }
  ```
- **Ответ:** `200` + обновлённая задача; недопустимый статус — `400` + `ErrInvalidStatus`; assignee не в команде — `400` + `ErrAssigneeNotInTeam`; нет прав — `403` / `404`.

---

### 1.5 История изменений задачи

- **Метод:** `GET /api/v1/tasks/{id}/history`
- **Сервис:** `TaskService.GetHistory(ctx, taskID, userID)`
- **Ответ:** `200` + массив записей истории; нет прав или задача не найдена — `403` / `404`.

---

## 2. Мягкое удаление и восстановление (расширение)

### 2.1 Мягкое удаление

- **Метод:** `DELETE /api/v1/tasks/{id}`
- **Сервис:** `TaskService.Delete(ctx, userID, taskID)`
- **Ответ:** `204` без тела; нет прав / задача не найдена — `403` / `404`.

---

### 2.2 Восстановление задачи

- **Метод:** `POST /api/v1/tasks/{id}/restore`
- **Сервис:** `TaskService.Restore(ctx, userID, taskID)`
- **Ответ:** `200` + восстановленная задача; нет прав / не найдена — `403` / `404`.

---

## 3. Сложные запросы / отчёты (обязательно по ТЗ)

### 3.1 Статистика по командам (п. а)

- **Назначение:** для каждой команды — название, кол-во участников, кол-во задач в статусе done за период (по ТЗ — последние 7 дней).
- **Метод:** `GET /api/v1/reports/team-task-stats?since=...`
  - Альтернатива: `GET /api/v1/teams/stats/tasks?since=...`
- **Сервис:** `TaskService.TeamTaskStats(ctx, userID, since)`
- **Query:** `since` — начало периода (ISO8601 или timestamp); по умолчанию — 7 дней назад.
- **Ответ:** `200` + массив `{ team_id, team_name, member_count, done_tasks_count }`.

---

### 3.2 Топ создателей задач по командам (п. б)

- **Назначение:** топ-N пользователей по количеству созданных задач в каждой команде за период (по ТЗ — за месяц).
- **Метод:** `GET /api/v1/reports/top-creators?since=...&limit=...`
  - Альтернатива: `GET /api/v1/teams/top-creators?since=...&limit=3`
- **Сервис:** `TaskService.TopCreatorsByTeam(ctx, userID, since, limit)`
- **Query:** `since`, `limit` (по ТЗ топ-3 → `limit=3` по умолчанию).
- **Ответ:** `200` + массив `{ team_id, user_id, rank, created_count }`.

---

### 3.3 Задачи с некорректным assignee (п. в)

- **Назначение:** задачи, у которых assignee не является членом команды задачи (валидация целостности).
- **Метод:** `GET /api/v1/reports/invalid-assignees`
  - Альтернатива: `GET /api/v1/tasks/invalid-assignees`
- **Сервис:** `TaskService.TasksWithInvalidAssignee(ctx, userID)`
- **Ответ:** `200` + массив задач (только по командам, где пользователь участник).

---

## 4. Сводная таблица ручек

| № | Метод | Путь | Сервис |
|---|--------|------|--------|
| 1 | POST | `/api/v1/tasks` | Create |
| 2 | GET | `/api/v1/tasks` | List |
| 3 | GET | `/api/v1/tasks/{id}` | GetByID |
| 4 | PUT | `/api/v1/tasks/{id}` | Update |
| 5 | DELETE | `/api/v1/tasks/{id}` | Delete |
| 6 | POST | `/api/v1/tasks/{id}/restore` | Restore |
| 7 | GET | `/api/v1/tasks/{id}/history` | GetHistory |
| 8 | GET | `/api/v1/reports/team-task-stats` | TeamTaskStats |
| 9 | GET | `/api/v1/reports/top-creators` | TopCreatorsByTeam |
| 10 | GET | `/api/v1/reports/invalid-assignees` | TasksWithInvalidAssignee |

---

## 5. Дополнительно по ТЗ

- **Кеширование Redis:** список задач команды (GET /api/v1/tasks с `team_id`) — TTL 5 мин.
- **Пагинация:** на уровне БД (LIMIT/OFFSET), параметры `limit`, `offset` в query.
- **Rate limiting:** 100 запросов/мин на пользователя (на все ручки).
- **Метрики Prometheus:** кол-во запросов, ошибок, время ответа по эндпоинтам.
