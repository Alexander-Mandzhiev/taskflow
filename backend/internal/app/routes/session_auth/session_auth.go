package session_auth

import (
	"context"

	"github.com/go-chi/chi/v5"

	account_v1 "mkk/internal/api/account/v1"
)

// Register регистрирует роуты account, защищённые session auth: logout.
// Whoami и другие эндпоинты для авторизованных пользователей добавляются сюда.
func Register(ctx context.Context, r chi.Router, api *account_v1.API) {
	r.Post("/logout", api.Logout)
}
