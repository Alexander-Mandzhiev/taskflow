package cache

import (
	"fmt"

	"github.com/google/uuid"
)

const keyPrefix = "refresh"

// Key возвращает ключ Redis для сессии по jti (типобезопасно — uuid.UUID).
func Key(jti uuid.UUID) string {
	return fmt.Sprintf("%s:%s", keyPrefix, jti.String())
}
