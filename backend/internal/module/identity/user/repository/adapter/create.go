package adapter

import (
	"context"

	"github.com/jmoiron/sqlx"

	"mkk/internal/module/identity/user/model"
	"mkk/internal/module/identity/user/repository/cache"
	"mkk/pkg/database/txmanager"
)

// Create создаёт пользователя в БД.
// Сохранение в кеш выполняется через post-commit hook в txmanager.
func (r *Repository) Create(ctx context.Context, tx *sqlx.Tx, input *model.UserInput, passwordHash string) (*model.User, error) {
	user, err := r.writer.Create(ctx, tx, input, passwordHash)
	if err != nil {
		return nil, err
	}
	registry := txmanager.GetHookRegistry(ctx)
	if registry != nil {
		id := user.ID.String()
		registry.Register(cache.Key(id), func(ctx context.Context) error {
			return r.cache.Set(ctx, id, user)
		})
	}
	return user, nil
}
