package session_auth

import (
	"context"

	"github.com/go-chi/chi/v5"

	account_v1 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/account/v1"
)

// Register регистрирует роуты account, защищённые JWT: logout.
func Register(ctx context.Context, r chi.Router, api *account_v1.API) {
	r.Post("/logout", api.Logout)
}
