package cache

import (
	"fmt"
	"time"
)

const keyPrefix = "user"

// TTL — время жизни записи пользователя в кеше (GetByID).
const TTL = 5 * time.Minute

// Key возвращает ключ Redis для кеша пользователя по id.
// Значение в кеше — JSON-сериализованный model.UserCache.
func Key(id string) string {
	return fmt.Sprintf("%s:%s", keyPrefix, id)
}
