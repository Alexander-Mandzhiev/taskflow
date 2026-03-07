package team

import (
	"context"

	"github.com/go-chi/chi/v5"

	team_v1 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/team/v1"
)

// Register регистрирует роуты команд. Вызывать только из группы с JWT и rate limit (registerAccountPrivateGroup).
// POST/GET /teams, GET /teams/{id}, POST /teams/{id}/invite.
func Register(ctx context.Context, r chi.Router, api *team_v1.API) {
	r.Post("/teams", api.Create)
	r.Get("/teams", api.List)
	r.Route("/teams/{id}", func(r chi.Router) {
		r.Get("/", api.GetByID)
		r.Post("/invite", api.Invite)
	})
}
