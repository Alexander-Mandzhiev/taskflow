package routes

import (
	"github.com/go-chi/chi/v5"

	account_v1 "mkk/internal/api/account/v1"
)

// RegisterAccountRoutes вешает хендлеры account API на роутер.
// Ожидается, что префикс уже задан (например r.Route("/api/account/v1", ...)).
func RegisterAccountRoutes(r chi.Router, api *account_v1.API) {
	r.Post("/register", api.Register)
	r.Post("/login", api.Login)
	r.Post("/logout", api.Logout)
}
