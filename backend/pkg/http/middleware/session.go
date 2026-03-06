package middleware

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/ctxkey"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
)

// SessionMiddleware извлекает session_id из cookie и устанавливает его в контекст
// Если cookie отсутствует или невалидна, запрос продолжается без session_id в контексте
// Это позволяет публичным эндпоинтам работать без сессии, а защищенным - проверять наличие сессии
func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(string(ctxkey.SessionID))
		if err == nil && cookie != nil && cookie.Value != "" {
			sessionID, err := uuid.Parse(cookie.Value)
			if err == nil {
				ctx := metadata.SetSessionIDUUID(r.Context(), sessionID)
				r = r.WithContext(ctx)
			}
		}

		next.ServeHTTP(w, r)
	})
}
