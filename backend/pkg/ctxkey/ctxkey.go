// Package ctxkey задаёт общие ключи контекста для использования в logger, metadata и других пакетах.
// Тип Key — свой тип (не string), чтобы избежать коллизий между пакетами (рекомендация Go context).
package ctxkey

// Key — тип ключа контекста. Используется как ключ в context.WithValue / ctx.Value.
type Key string

const (
	// TraceID — глобальный идентификатор трассировки запроса (логирование, трейсинг).
	TraceID Key = "trace_id"
	// RequestID — уникальный идентификатор HTTP/gRPC запроса.
	RequestID Key = "request_id"
	// UserID — идентификатор пользователя (сессия, логирование).
	UserID Key = "user_id"
	// SessionID — идентификатор сессии (аутентификация).
	SessionID Key = "session_id"
)
