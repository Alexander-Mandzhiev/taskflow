package adapter

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
	usercache "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository/cache/user"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/txmanager"
)

// Update обновляет пользователя в БД.
// Обновление кеша выполняется через post-commit hook в txmanager.
func (r *Adapter) Update(ctx context.Context, tx *sqlx.Tx, id string, input *model.UserInput) (*model.User, error) {
	user, err := r.writer.Update(ctx, tx, id, input)
	if err != nil {
		return nil, err
	}
	registry := txmanager.GetHookRegistry(ctx)
	if registry != nil {
		registry.Register(usercache.Key(id), func(ctx context.Context) error {
			return r.cache.Set(ctx, id, user)
		})
	}
	return user, nil
}
