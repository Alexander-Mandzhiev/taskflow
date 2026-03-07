package task

import (
	"context"

	"github.com/go-chi/chi/v5"

	task_v1 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/task/v1"
)

// Register регистрирует роуты задач, отчётов и комментариев. Вызывать только из группы с JWT и rate limit.
// POST/GET /tasks, GET/PUT /tasks/{id}, GET /tasks/{id}/history, GET /reports/..., GET/POST /tasks/{id}/comments.
func Register(ctx context.Context, r chi.Router, api *task_v1.API) {
	r.Post("/tasks", api.Create)
	r.Get("/tasks", api.List)
	r.Route("/tasks/{id}", func(r chi.Router) {
		r.Get("/", api.GetByID)
		r.Put("/", api.Update)
		r.Get("/history", api.GetHistory)
		r.Get("/comments", api.ListComments)
		r.Post("/comments", api.CreateComment)
	})
	r.Route("/reports", func(r chi.Router) {
		r.Get("/team-stats", api.ReportTeamStats)
		r.Get("/top-creators", api.ReportTopCreators)
		r.Get("/invalid-assignees", api.ReportInvalidAssignees)
	})
}
