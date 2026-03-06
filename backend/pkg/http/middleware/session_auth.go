package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	accountmodel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/ctxkey"
	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
)

// SessionWhoamiService определяет интерфейс сервиса для проверки сессии (auth по session_id).
type SessionWhoamiService interface {
	Whoami(ctx context.Context, sessionID uuid.UUID) (uuid.UUID, error)
}

// SessionAuthMiddleware выполняет проверку аутентификации на основе session_id из контекста
// и сервиса сессий (Whoami). При успешной проверке записывает user_id в контекст.
// Используется при auth по сессиям; для JWT используйте JWTAuthMiddleware.
type SessionAuthMiddleware struct {
	sessionService SessionWhoamiService
	isSecure       bool
	cookieDomain   string
}

// NewSessionAuthMiddleware создаёт middleware для проверки сессии.
func NewSessionAuthMiddleware(sessionService SessionWhoamiService, isSecure bool, cookieDomain string) *SessionAuthMiddleware {
	return &SessionAuthMiddleware{
		sessionService: sessionService,
		isSecure:       isSecure,
		cookieDomain:   cookieDomain,
	}
}

// Handle проверяет, что запрос аутентифицирован сессией.
func (m *SessionAuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := metadata.SessionID(r.Context())
		if err != nil || sessionID == uuid.Nil {
			pkghttp.DeleteCookie(w, string(ctxkey.SessionID), m.isSecure, m.cookieDomain)
			writeAuthResponse(r.Context(), w, false, "Требуется аутентификация")
			return
		}

		userID, err := m.sessionService.Whoami(r.Context(), sessionID)
		if err != nil {
			if errors.Is(err, accountmodel.ErrSessionNotFound) {
				pkghttp.DeleteCookie(w, string(ctxkey.SessionID), m.isSecure, m.cookieDomain)
				writeAuthResponse(r.Context(), w, false, "Сессия не найдена или истекла")
				return
			}
			pkghttp.WriteJSON(r.Context(), w, http.StatusInternalServerError, pkghttp.ErrorBody{
				Code:    http.StatusInternalServerError,
				Message: "Внутренняя ошибка сервера",
			})
			return
		}

		ctx := metadata.SetUserID(r.Context(), userID.String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
