package middleware

import (
	"errors"
	"net/http"
	"strings"

	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/jwt"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
)

// JWTAuthMiddleware выполняет проверку аутентификации по JWT:
// извлекает токен из заголовка Authorization: Bearer <token> или из cookie (имя из конфига),
// валидирует, записывает user_id в контекст. При ошибке — 401 и удаление cookie при наличии.
type JWTAuthMiddleware struct {
	secret       string
	isSecure     bool
	cookieDomain string
	cookieName   string // имя cookie для access-токена (например "access_token")
}

// NewJWTAuthMiddleware создаёт middleware для проверки JWT.
func NewJWTAuthMiddleware(secret string, isSecure bool, cookieDomain, accessTokenCookieName string) *JWTAuthMiddleware {
	cookieName := accessTokenCookieName
	if cookieName == "" {
		cookieName = pkghttp.DefaultAccessTokenCookieName
	}
	return &JWTAuthMiddleware{
		secret:       secret,
		isSecure:     isSecure,
		cookieDomain: cookieDomain,
		cookieName:   cookieName,
	}
}

// Handle проверяет JWT и устанавливает user_id в контекст.
func (m *JWTAuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := m.extractToken(r)
		if token == "" {
			m.deleteTokenCookie(w)
			writeAuthResponse(r.Context(), w, "Требуется аутентификация")
			return
		}

		claims, err := jwt.ValidateToken(token, m.secret)
		if err != nil {
			m.deleteTokenCookie(w)
			msg := "Неверный токен"
			if errors.Is(err, jwt.ErrExpiredToken) {
				msg = "Токен истёк"
			}
			writeAuthResponse(r.Context(), w, msg)
			return
		}

		ctx := metadata.SetUserID(r.Context(), claims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// extractToken извлекает JWT: сначала из заголовка Authorization: Bearer <token>, затем из cookie.
func (m *JWTAuthMiddleware) extractToken(r *http.Request) string {
	if h := r.Header.Get("Authorization"); h != "" {
		if parts := strings.SplitN(h, " ", 2); len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") && parts[1] != "" {
			return parts[1]
		}
	}
	if c, err := r.Cookie(m.cookieName); err == nil && c.Value != "" {
		return c.Value
	}
	return ""
}

func (m *JWTAuthMiddleware) deleteTokenCookie(w http.ResponseWriter) {
	pkghttp.DeleteCookie(w, m.cookieName, m.isSecure, m.cookieDomain)
}
