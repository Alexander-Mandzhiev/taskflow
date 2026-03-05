package service

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	usermodel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

// Register создаёт пользователя: хеш пароля — до транзакции; в транзакции проверка email и создание записи.
// Дубликат email — usermodel.ErrEmailDuplicate (Create тоже вернёт его при гонке).
func (s *accountService) Register(ctx context.Context, email, password, name string) error {
	// 1. Хешируем ДО транзакции. CPU работает, БД отдыхает.
	hash, err := s.hasher.Hash(password)
	if err != nil {
		logger.Error(ctx, "Register: hash password failed", zap.Error(err))
		return err
	}

	// 2. Быстрая транзакция
	err = s.txManager.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// Проверка для UX: избежать лишнего Insert, если email уже занят
		existing, errGet := s.userRepo.GetByEmail(ctx, tx, email)
		if errGet != nil && !errors.Is(errGet, usermodel.ErrUserNotFound) {
			return errGet
		}
		if existing != nil {
			return usermodel.ErrEmailDuplicate
		}

		input := &usermodel.UserInput{Email: email, Name: name}
		// Create сам проверит дубликат через БД и вернёт ErrEmailDuplicate,
		// если кто-то зарегистрировался за миллисекунды после GetByEmail
		_, errCreate := s.userRepo.Create(ctx, tx, input, hash)
		return errCreate
	})

	if err != nil {
		// Логируем только системные ошибки. Дубликат email — ожидаемая бизнес-ошибка.
		if !errors.Is(err, usermodel.ErrEmailDuplicate) {
			logger.Error(ctx, "Register: process failed", zap.Error(err))
		}
		return err
	}

	return nil
}
