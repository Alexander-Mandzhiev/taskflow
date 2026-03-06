package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/ctxkey"
	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/jwt"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
)

// JWTAuthMiddleware выполняет проверку аутентификации по JWT:
// извлекает токен из заголовка Authorization: Bearer <token> или из cookie access_token,
// валидирует, записывает user_id в контекст. При ошибке — 401 и удаление cookie при наличии.
type JWTAuthMiddleware struct {
	secret       string
	isSecure     bool
	cookieDomain string
}

// NewJWTAuthMiddleware создаёт middleware для проверки JWT.
func NewJWTAuthMiddleware(secret string, isSecure bool, cookieDomain string) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{
		secret:       secret,
		isSecure:     isSecure,
		cookieDomain: cookieDomain,
	}
}

// Handle проверяет JWT и устанавливает user_id в контекст.
func (m *JWTAuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			m.deleteTokenCookie(w)
			writeAuthResponse(r.Context(), w, false, "Требуется аутентификация")
			return
		}

		claims, err := jwt.ValidateToken(token, m.secret)
		if err != nil {
			m.deleteTokenCookie(w)
			msg := "Невалидный или истёкший токен"
			if errors.Is(err, jwt.ErrExpiredToken) {
				msg = "Токен истёк"
			}
			writeAuthResponse(r.Context(), w, false, msg)
			return
		}

		ctx := metadata.SetUserID(r.Context(), claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractToken(r *http.Request) string {
	if h := r.Header.Get("Authorization"); h != "" {
		if parts := strings.SplitN(h, " ", 2); len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") && parts[1] != "" {
			return parts[1]
		}
	}
	if c, err := r.Cookie(string(ctxkey.AccessToken)); err == nil && c != nil && c.Value != "" {
		return c.Value
	}
	return ""
}

func (m *JWTAuthMiddleware) deleteTokenCookie(w http.ResponseWriter) {
	pkghttp.DeleteCookie(w, string(ctxkey.AccessToken), m.isSecure, m.cookieDomain)
}
