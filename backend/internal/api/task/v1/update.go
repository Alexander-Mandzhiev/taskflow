package task_v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/task/v1/dto"
	taskconverter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/converter"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/util"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/validation"
)

// Update обновляет задачу. PUT /api/v1/tasks/{id}. Body: UpdateTaskRequest.
func (api *API) Update(w http.ResponseWriter, r *http.Request) {
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

	var req dto.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mapError(w, r, err)
		return
	}

	if err := validation.Validator.Struct(req); err != nil {
		mapError(w, r, err)
		return
	}

	input, err := taskconverter.UpdateTaskRequestToDomain(req)
	if err != nil {
		mapError(w, r, err)
		return
	}

	task, err := api.taskService.Update(r.Context(), userID, taskID, input)
	if err != nil {
		mapError(w, r, err)
		return
	}

	pkghttp.WriteJSON(r.Context(), w, http.StatusOK, taskconverter.TaskToResponse(task))
}
