package adapter

import (
	"context"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

// GetByID возвращает пользователя по ID.
// При tx != nil — чтение из БД в транзакции (валидация внутри мутаций).
// При tx == nil — кеш, при промахе или ошибке Redis — fallback на БД.
// PasswordHash в кеше не хранится: для аутентификации используйте GetByEmail.
func (r *Adapter) GetByID(ctx context.Context, tx *sqlx.Tx, id string) (*model.User, error) {
	if tx != nil {
		return r.reader.GetByID(ctx, tx, id)
	}

	// 1. Пытаемся получить из кеша
	user, err := r.cache.Get(ctx, id)
	if err != nil {
		logger.Warn(ctx, "Cache get failed, falling back to DB", zap.Error(err))
	}
	if user != nil {
		return user, nil
	}

	// 2. Fallback на БД
	user, err = r.reader.GetByID(ctx, nil, id)
	if err != nil {
		return nil, err
	}

	// 3. Кешируем результат
	if err := r.cache.Set(ctx, id, user); err != nil {
		logger.Warn(ctx, "Cache set after get by id failed", zap.Error(err))
	}

	return user, nil
}
