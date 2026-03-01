package cache

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Delete удаляет сессию из кеша (logout).
func (r *repository) Delete(ctx context.Context, sessionID uuid.UUID) error {
	key := Key(sessionID)
	if err := r.redis.Del(ctx, key); err != nil {
		return fmt.Errorf("session cache delete: %w", err)
	}
	return nil
}
