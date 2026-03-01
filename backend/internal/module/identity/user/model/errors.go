package model

import "errors"

// Ошибки домена пользователя. Сервисный и API-слой проверяют через errors.Is.
var (
	// ErrUserNotFound — пользователь не найден (по id или email).
	ErrUserNotFound = errors.New("user not found")

	// ErrNilInput — передан nil вместо *UserInput в Create/Update.
	ErrNilInput = errors.New("user input is nil")

	// ErrEmailDuplicate — email уже занят (нарушение уникальности при создании или обновлении).
	ErrEmailDuplicate = errors.New("email already exists")

	// ErrInvalidCredentials — неверный email или пароль (логин).
	ErrInvalidCredentials = errors.New("invalid credentials")
)
