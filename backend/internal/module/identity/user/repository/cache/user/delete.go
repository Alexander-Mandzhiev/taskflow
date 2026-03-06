package cache

import (
	"context"
	"fmt"
)

// Delete удаляет запись пользователя из кеша (инвалидация при Update/Delete).
func (r *repository) Delete(ctx context.Context, id string) error {
	key := Key(id)
	if err := r.redis.Del(ctx, key); err != nil {
		return fmt.Errorf("cache delete: %w", err)
	}
	return nil
}
