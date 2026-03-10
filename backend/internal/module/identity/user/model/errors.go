package model

import "errors"

// Ошибки домена пользователя. Сервисный и API-слой проверяют через errors.Is.
var (
	// ErrUserNotFound — пользователь не найден (по id или email).
	ErrUserNotFound = errors.New("user not found")

	// ErrNilInput — передан нулевой UserInput в Create/Update (пустые Email и Name).
	ErrNilInput = errors.New("user input is nil")

	// ErrEmailDuplicate — email уже занят (нарушение уникальности при создании или обновлении).
	ErrEmailDuplicate = errors.New("email already exists")

	// ErrInvalidCredentials — неверный email или пароль (логин).
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrInvalidInput — невалидные данные (например NOT NULL нарушение на уровне БД).
	ErrInvalidInput = errors.New("invalid input")

	// ErrTemporaryFailure — временная ошибка (deadlock, lock wait); можно повторить запрос.
	ErrTemporaryFailure = errors.New("temporary failure")

	// ErrInternal — внутренняя/неизвестная ошибка (в т.ч. прочая ошибка БД).
	ErrInternal = errors.New("internal error")
)
