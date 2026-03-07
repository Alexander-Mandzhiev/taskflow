package model

import "github.com/google/uuid"

// ValidateTaskInput проверяет допустимость полей input. Если задан статус — он должен быть todo, in_progress или done.
// Для использования в сервисе при Update (при Create валидация статуса выполняется в конвертере). Возвращает ErrInvalidStatus при недопустимом статусе.
func ValidateTaskInput(input *TaskInput) error {
	if input == nil {
		return nil
	}
	if input.Status != "" && !IsValidTaskStatus(input.Status) {
		return ErrInvalidStatus
	}
	return nil
}

// TaskInput — единая модель для создания и обновления задачи (тело запроса).
// Create: teamID передаётся в сигнатуре. Update: taskID передаётся в сигнатуре.
// AssigneeID — опционально (nil = не назначен / снять исполнителя). Status по умолчанию todo — в сервисе при необходимости.
type TaskInput struct {
	Title       string
	Description string
	Status      string
	AssigneeID  *uuid.UUID
}

// TaskListFilter — фильтр и пагинация для списка задач (GET /api/v1/tasks?team_id=...&status=...&assignee_id=...&limit=...&offset=...).
// Валидация (limit > 0 и т.д.) — в сервисе или API/DTO.
type TaskListFilter struct {
	TeamID     *uuid.UUID
	Status     *string
	AssigneeID *uuid.UUID
	Limit      int // обязателен с фронта, проверка limit > 0 в сервисе/API
	Offset     int
}
