package cache

import (
	"fmt"

	"github.com/google/uuid"
)

const keyPrefix = "session"

// Key возвращает ключ Redis для сессии по sessionID.
// Значение в кеше — JSON resources.SessionCache.
func Key(sessionID uuid.UUID) string {
	return fmt.Sprintf("%s:%s", keyPrefix, sessionID.String())
}
