package account_v1

import (
	"time"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/service"
)

// API — HTTP-слой account (логин, логаут). Токены приходят из сервиса; API выставляет/удаляет cookie.
// Logout: читает refresh из cookie (сырая строка), передаёт в сервис; валидация и удаление сессии — в сервисе.
type API struct {
	accountService         service.AccountService
	accessTokenCookieName  string
	accessTTL              time.Duration
	refreshTokenCookieName string
	refreshTTL             time.Duration
	isSecure               bool
	cookieDomain           string
}

// NewAPI создаёт API.
func NewAPI(
	accountService service.AccountService,
	accessTokenCookieName string,
	accessTTL time.Duration,
	refreshTokenCookieName string,
	refreshTTL time.Duration,
	isSecure bool,
	cookieDomain string,
) *API {
	return &API{
		accountService:         accountService,
		accessTokenCookieName:  accessTokenCookieName,
		accessTTL:              accessTTL,
		refreshTokenCookieName: refreshTokenCookieName,
		refreshTTL:             refreshTTL,
		isSecure:               isSecure,
		cookieDomain:           cookieDomain,
	}
}
