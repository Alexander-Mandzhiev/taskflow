package cache

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Delete удаляет сессию из кеша (logout).
func (r *repository) Delete(ctx context.Context, jti uuid.UUID) error {
	key := Key(jti)
	if err := r.redis.Del(ctx, key); err != nil {
		return fmt.Errorf("session cache delete: %w", err)
	}
	return nil
}
