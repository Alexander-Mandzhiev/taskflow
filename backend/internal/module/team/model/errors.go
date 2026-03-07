package model

import "errors"

// Ошибки домена команды. Сервисный и API-слой проверяют через errors.Is.
var (
	// ErrTeamNotFound — команда не найдена по id.
	ErrTeamNotFound = errors.New("team not found")

	// ErrMemberNotFound — участник не найден (пара team_id, user_id).
	ErrMemberNotFound = errors.New("member not found")

	// ErrAlreadyMember — пользователь уже является участником команды (дубликат при invite).
	ErrAlreadyMember = errors.New("user is already a member of the team")

	// ErrForbidden — недостаточно прав (например, приглашать может только owner/admin).
	ErrForbidden = errors.New("forbidden")

	// ErrNilInput — передан nil вместо *TeamInput в Create.
	ErrNilInput = errors.New("team input is nil")

	// ErrTemporaryFailure — временная ошибка (deadlock, lock wait); можно повторить запрос.
	ErrTemporaryFailure = errors.New("temporary failure")

	// ErrInternal — внутренняя/неизвестная ошибка (в т.ч. прочая ошибка БД).
	ErrInternal = errors.New("internal error")

	// ErrUserNotFound — пользователь с указанным email не найден (при invite по email).
	ErrUserNotFound = errors.New("user not found")

	// ErrInvalidID — некорректный идентификатор в URL (не UUID или пустой).
	ErrInvalidID = errors.New("invalid id")
)
