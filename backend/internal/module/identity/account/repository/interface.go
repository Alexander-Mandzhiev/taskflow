package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/model"
)

// SessionCacheRepository — хранилище сессий (ключ session:{session_id}, значение — model.Session, TTL на ключе).
// Единственный слой персистентности для сессий: без таблицы сессий в БД.
// Get при отсутствии/истечении сессии возвращает model.ErrSessionNotFound.
type SessionCacheRepository interface {
	// Set создаёт или обновляет сессию. session — user_id, created_at, device_type, user_agent, ip.
	Set(ctx context.Context, sessionID uuid.UUID, session *model.Session, ttl time.Duration) error
	// Get возвращает данные сессии по sessionID. Для Whoami достаточно session.UserID.
	Get(ctx context.Context, sessionID uuid.UUID) (*model.Session, error)
	// Delete удаляет сессию (logout).
	Delete(ctx context.Context, sessionID uuid.UUID) error
}
