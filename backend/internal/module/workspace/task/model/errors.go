package model

import "errors"

// Ошибки домена задач. Сервисный и API-слой проверяют через errors.Is.
var (
	// ErrTaskNotFound — задача не найдена по id.
	ErrTaskNotFound = errors.New("task not found")

	// ErrForbidden — недостаточно прав (пользователь не в команде задачи или нет прав на действие).
	ErrForbidden = errors.New("forbidden")

	// ErrNilInput — передан nil вместо входной структуры (Create/Update).
	ErrNilInput = errors.New("task input is nil")

	// ErrInvalidStatus — недопустимый статус задачи (допустимы todo, in_progress, done).
	ErrInvalidStatus = errors.New("invalid task status")

	// ErrAssigneeNotInTeam — assignee не является участником команды задачи.
	ErrAssigneeNotInTeam = errors.New("assignee is not a member of the task team")

	// ErrTemporaryFailure — временная ошибка (deadlock, lock wait); можно повторить запрос.
	ErrTemporaryFailure = errors.New("temporary failure")

	// ErrInternal — внутренняя/неизвестная ошибка.
	ErrInternal = errors.New("internal error")
)
