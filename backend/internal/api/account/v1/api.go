package account_v1

import (
	"time"

	"mkk/internal/module/account"
)

// API реализует HTTP‑хендлеры для account: регистрация, логин, логаут.
type API struct {
	accountService account.Service
	sessionTTL     time.Duration
	isSecure       bool
	cookieDomain   string
}

// NewAPI создаёт API account v1 с сервисом account и настройками cookie.
func NewAPI(
	accountService account.Service,
	sessionTTL time.Duration,
	isSecure bool,
	cookieDomain string,
) *API {
	return &API{
		accountService: accountService,
		sessionTTL:     sessionTTL,
		isSecure:       isSecure,
		cookieDomain:   cookieDomain,
	}
}
