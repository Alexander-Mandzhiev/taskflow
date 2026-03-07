package adapter

import (
	"context"

	"github.com/jmoiron/sqlx"

	usercache "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository/cache/user"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/txmanager"
)

// Delete удаляет пользователя в БД (мягкое удаление).
// Инвалидация кеша выполняется через post-commit hook в txmanager.
func (r *Adapter) Delete(ctx context.Context, tx *sqlx.Tx, id string) error {
	if err := r.writer.Delete(ctx, tx, id); err != nil {
		return err
	}
	registry := txmanager.GetHookRegistry(ctx)
	if registry != nil {
		userID := id
		registry.Register(usercache.Key(userID), func(ctx context.Context) error {
			return r.cache.Delete(ctx, userID)
		})
	}
	return nil
}
