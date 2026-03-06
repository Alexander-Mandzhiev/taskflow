package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/repository/converter"
)

// Set создаёт или обновляет сессию в кеше по jti с заданным TTL.
func (r *repository) Set(ctx context.Context, jti uuid.UUID, session *model.Session, ttl time.Duration) error {
	if session == nil {
		return nil
	}
	c := converter.ToCache(session)
	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("session cache marshal: %w", err)
	}
	if err := r.redis.Set(ctx, Key(jti), data, ttl); err != nil {
		return fmt.Errorf("session cache set: %w", err)
	}
	return nil
}
