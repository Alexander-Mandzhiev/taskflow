package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/jwt"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
)

// RefreshTokenRefresher — контракт для тихого рефреша: по валидному refresh выдать новый access и userID.
// Реализуется, например, AccountService.
type RefreshTokenRefresher interface {
	Refresh(ctx context.Context, refreshToken, userAgent, ip string) (accessToken string, userID uuid.UUID, err error)
}

// JWTAuthMiddleware выполняет проверку аутентификации по JWT:
// извлекает access из заголовка Authorization или cookie, валидирует, записывает user_id в контекст.
// При невалидном/истёкшем access выполняет тихий рефреш (если передан refresher): по refresh из cookie
// выдаёт новый access и продолжает запрос; при ошибке рефреша — 401.
type JWTAuthMiddleware struct {
	secret            string
	isSecure          bool
	cookieDomain      string
	cookieName        string // имя cookie для access-токена
	refresher         RefreshTokenRefresher
	refreshCookieName string // имя cookie для refresh (тихий рефреш)
	accessTTLSeconds  int    // max-age для новой access-cookie при рефреше
}

// NewJWTAuthMiddleware создаёт middleware для проверки JWT.
// refresher — опционально: при nil тихий рефреш не выполняется; иначе при ошибке access пробуем refresh.
// refreshCookieName и accessTTL используются только при ненулевом refresher.
func NewJWTAuthMiddleware(
	secret string,
	isSecure bool,
	cookieDomain string,
	accessTokenCookieName string,
	refresher RefreshTokenRefresher,
	refreshCookieName string,
	accessTTL time.Duration,
) *JWTAuthMiddleware {
	cookieName := accessTokenCookieName
	if cookieName == "" {
		cookieName = pkghttp.DefaultAccessTokenCookieName
	}
	accessTTLSeconds := 0
	if accessTTL > 0 {
		accessTTLSeconds = int(accessTTL.Seconds())
	}
	return &JWTAuthMiddleware{
		secret:            secret,
		isSecure:          isSecure,
		cookieDomain:      cookieDomain,
		cookieName:        cookieName,
		refresher:         refresher,
		refreshCookieName: refreshCookieName,
		accessTTLSeconds:  accessTTLSeconds,
	}
}

// Handle проверяет JWT и устанавливает user_id в контекст. При невалидном access пробует тихий рефреш (если refresher задан).
func (m *JWTAuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := m.extractToken(r)
		if token != "" {
			claims, err := jwt.ValidateToken(token, m.secret)
			if err == nil {
				ctx := metadata.SetUserID(r.Context(), claims.Subject)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		// Access отсутствует или невалиден — пробуем тихий рефреш
		if m.refresher != nil && m.refreshCookieName != "" && m.accessTTLSeconds > 0 {
			if refreshToken := m.extractRefreshToken(r); refreshToken != "" {
				accessToken, userID, err := m.refresher.Refresh(r.Context(), refreshToken, r.UserAgent(), pkghttp.ClientIP(r))
				if err == nil {
					pkghttp.SetCookie(w, m.cookieName, accessToken, m.accessTTLSeconds, m.isSecure, m.cookieDomain, false)
					ctx := metadata.SetUserID(r.Context(), userID.String())
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}
		}

		m.deleteTokenCookie(w)
		msg := "Требуется аутентификация"
		if token != "" {
			msg = "Токен истёк или недействителен"
		}
		writeAuthResponse(r.Context(), w, msg)
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

func (m *JWTAuthMiddleware) extractRefreshToken(r *http.Request) string {
	if c, err := r.Cookie(m.refreshCookieName); err == nil && c.Value != "" {
		return c.Value
	}
	return ""
}
