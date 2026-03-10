package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/model"
)

// SessionCacheRepository — хранилище сессий (ключ — jti как uuid.UUID).
// Значение — model.Session (метаданные), TTL на ключе. Токены в Redis не храним.
// Get при отсутствии/истечении возвращает (model.Session{}, model.ErrSessionNotFound).
type SessionCacheRepository interface {
	// Set создаёт или обновляет сессию. session — user_id, created_at, device_type, user_agent, ip.
	Set(ctx context.Context, jti uuid.UUID, session model.Session, ttl time.Duration) error
	// Get возвращает данные сессии по jti (например для проверки при refresh).
	Get(ctx context.Context, jti uuid.UUID) (model.Session, error)
	// Delete удаляет сессию (logout).
	Delete(ctx context.Context, jti uuid.UUID) error
}
