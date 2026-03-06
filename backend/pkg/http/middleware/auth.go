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

// SessionWhoamiService определяет интерфейс сервиса для проверки сессии.

type SessionWhoamiService interface {
	Whoami(ctx context.Context, sessionID uuid.UUID) (uuid.UUID, error)
}

// AuthMiddleware выполняет проверку аутентификации на основе session_id из контекста

// и сервиса сессий (Whoami). При успешной проверке записывает user_id в контекст.

type AuthMiddleware struct {
	sessionService SessionWhoamiService

	isSecure bool

	cookieDomain string
}

// NewAuthMiddleware создаёт новый экземпляр AuthMiddleware.

func NewAuthMiddleware(sessionService SessionWhoamiService, isSecure bool, cookieDomain string) *AuthMiddleware {
	return &AuthMiddleware{
		sessionService: sessionService,

		isSecure: isSecure,

		cookieDomain: cookieDomain,
	}
}

// Handle проверяет, что запрос аутентифицирован сессией.

// 1. Достаёт session_id из контекста (его туда положил SessionMiddleware)

// 2. Вызывает Whoami для проверки сессии

// 3. Кладёт user_id в контекст

// 4. При ошибке Whoami: если сессия не найдена/истекла (ErrSessionNotFound) — удаляет куку и 401;
//    при внутренней ошибке (таймаут БД и т.п.) — 500 без удаления cookie, чтобы не разлогинивать пользователя.

func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
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
			// Внутренняя ошибка (таймаут БД, недоступность хранилища и т.п.) — не удаляем cookie.
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

// AuthResponse представляет ответ при невалидной сессии.

type AuthResponse struct {
	Authenticated bool   `json:"authenticated"`
	Message       string `json:"message,omitempty"`
}

func writeAuthResponse(ctx context.Context, w http.ResponseWriter, authenticated bool, message string) {
	code := http.StatusOK
	if !authenticated {
		code = http.StatusUnauthorized
	}

	pkghttp.WriteJSON(ctx, w, code, AuthResponse{
		Authenticated: authenticated,
		Message:       message,
	})
}
