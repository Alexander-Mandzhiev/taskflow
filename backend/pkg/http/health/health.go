package healthhttp

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
)

// Response структура ответа для health check
type Response struct {
	Status string `json:"status"`
}

// Handler возвращает простой JSON-ответ статуса для healthcheck.
func Handler(w http.ResponseWriter, r *http.Request) {
	pkghttp.WriteJSON(w, http.StatusOK, Response{Status: "ok"})
}

// LiveHandler — эндпоинт liveness (приложение запущено).
func LiveHandler(w http.ResponseWriter, r *http.Request) {
	Handler(w, r)
}

// ReadyHandler — эндпоинт readiness (приложение готово обслуживать трафик).
func ReadyHandler(w http.ResponseWriter, r *http.Request) {
	Handler(w, r)
}

// StartHandler — эндпоинт startup (инициализация завершена).
func StartHandler(w http.ResponseWriter, r *http.Request) {
	Handler(w, r)
}

// RegisterRoutes регистрирует стандартные health check маршруты на роутере.
// Поддерживает стандартные endpoints для Kubernetes и других оркестраторов:
// - /health - простой health check (возвращает JSON {"status": "ok"})
// - /healthz - Kubernetes liveness probe
// - /live - Kubernetes liveness probe (альтернативный)
// - /ready - Kubernetes readiness probe
// - /start - Kubernetes startup probe
func RegisterRoutes(r chi.Router) {
	r.Get("/health", Handler)
	r.Get("/healthz", Handler)
	r.Get("/live", LiveHandler)
	r.Get("/ready", ReadyHandler)
	r.Get("/start", StartHandler)
}
