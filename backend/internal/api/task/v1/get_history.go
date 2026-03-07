package task_v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	taskconverter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/converter"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/util"
)

// GetHistory возвращает историю изменений задачи. GET /api/v1/tasks/{id}/history.
func (api *API) GetHistory(w http.ResponseWriter, r *http.Request) {
	userID, err := metadata.UserID(r.Context())
	if err != nil {
		mapError(w, r, err)
		return
	}
	taskID, err := util.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		mapError(w, r, model.ErrTaskNotFound)
		return
	}

	entries, err := api.taskService.GetHistory(r.Context(), taskID, userID)
	if err != nil {
		mapError(w, r, err)
		return
	}

	pkghttp.WriteJSON(r.Context(), w, http.StatusOK, taskconverter.TaskHistoryToResponse(taskID.String(), entries))
}
