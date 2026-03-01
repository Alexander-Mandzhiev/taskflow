package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"mkk/internal/module/identity/account/model"
	"mkk/internal/module/identity/account/repository/converter"
)

// Set создаёт или обновляет сессию в кеше с заданным TTL.
func (r *repository) Set(ctx context.Context, sessionID uuid.UUID, session *model.Session, ttl time.Duration) error {
	if session == nil {
		return nil
	}
	key := Key(sessionID)
	c := converter.ToCache(session)
	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("session cache marshal: %w", err)
	}
	if err := r.redis.Set(ctx, key, data, ttl); err != nil {
		return fmt.Errorf("session cache set: %w", err)
	}
	return nil
}
