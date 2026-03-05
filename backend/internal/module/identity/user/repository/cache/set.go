package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository/converter"
)

// Set сохраняет пользователя в кеш по id с TTL.
func (r *repository) Set(ctx context.Context, id string, user *model.User) error {
	if user == nil {
		return nil
	}
	key := Key(id)
	c := converter.ToCache(*user)
	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("cache marshal: %w", err)
	}
	if err := r.redis.Set(ctx, key, data, TTL); err != nil {
		return fmt.Errorf("cache set: %w", err)
	}
	return nil
}
