package adapter

import (
	"context"

	"github.com/jmoiron/sqlx"

	"mkk/internal/module/identity/user/repository/cache"
	"mkk/pkg/database/txmanager"
)

// Delete удаляет пользователя в БД (мягкое удаление).
// Инвалидация кеша выполняется через post-commit hook в txmanager.
func (r *Repository) Delete(ctx context.Context, tx *sqlx.Tx, id string) error {
	if err := r.writer.Delete(ctx, tx, id); err != nil {
		return err
	}
	registry := txmanager.GetHookRegistry(ctx)
	if registry != nil {
		userID := id
		registry.Register(cache.Key(userID), func(ctx context.Context) error {
			return r.cache.Delete(ctx, userID)
		})
	}
	return nil
}
