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

	// ErrPaginationRequired — для List обязательна пагинация с фронта: limit > 0.
	ErrPaginationRequired = errors.New("pagination required: limit must be positive")

	// ErrTeamIDRequired — для списка задач обязателен query-параметр team_id.
	ErrTeamIDRequired = errors.New("team_id is required")

	// ErrInvalidAssigneeID — передан невалидный assignee_id в query (ожидается UUID).
	ErrInvalidAssigneeID = errors.New("invalid assignee_id parameter")

	// ErrInvalidLimit — limit должен быть положительным (например, в отчётах).
	ErrInvalidLimit = errors.New("limit must be positive")

	// ErrTxRequired — мутация вызвана без транзакции (tx == nil). Writer-репозитории требуют tx.
	ErrTxRequired = errors.New("transaction required")

	// ErrCommentNotImplemented — сервис комментариев пока не реализован (заглушка).
	ErrCommentNotImplemented = errors.New("comment service not implemented")

	// ErrCommentNotFound — комментарий не найден по id.
	ErrCommentNotFound = errors.New("comment not found")
)
