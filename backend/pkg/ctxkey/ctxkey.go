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
	// SessionKey — ключ сессии в контексте (uuid: jti в JWT flow или session_id в legacy). Используется для logout.
	SessionKey Key = "session_key"
	// AccessToken — имя cookie/заголовка для JWT access токена.
	AccessToken Key = "access_token"
	// RefreshToken — имя httpOnly cookie для JWT refresh токена (обновление access без перелогина).
	RefreshToken Key = "refresh_token"
)
