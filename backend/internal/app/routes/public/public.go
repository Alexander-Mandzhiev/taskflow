package public

import (
	"context"

	"github.com/go-chi/chi/v5"

	account_v1 "mkk/internal/api/account/v1"
)

// Register регистрирует публичные роуты account (без проверки сессии): register, login.
func Register(ctx context.Context, r chi.Router, api *account_v1.API) {
	r.Post("/register", api.Register)
	r.Post("/login", api.Login)
}
