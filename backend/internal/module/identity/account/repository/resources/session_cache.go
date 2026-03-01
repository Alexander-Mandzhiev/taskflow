package resources

import "time"

// SessionCache — модель сессии для хранения в кеше (Redis, ключ session:{id}).
// Отдельная структура от model.Session позволяет менять формат хранения без изменения домена.
type SessionCache struct {
	UserID     string    `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	DeviceType string    `json:"device_type,omitempty"`
	UserAgent  string    `json:"user_agent,omitempty"`
	IP         string    `json:"ip,omitempty"`
}
