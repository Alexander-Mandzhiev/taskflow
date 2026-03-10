package repository

import (
	"context"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
)

// UserCacheRepository предоставляет методы для кеша пользователей (по id).
// Get — чтение из кеша (found == false при промахе), Set — запись после чтения из БД, Delete — инвалидация при записи в БД.
type UserCacheRepository interface {
	Get(ctx context.Context, id string) (model.User, bool, error)
	Set(ctx context.Context, id string, user model.User) error
	Delete(ctx context.Context, id string) error
}
