package task_v1

import (
	"encoding/json"
	"net/http"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/task/v1/dto"
	taskconverter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/converter"
	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/validation"
)

// Create создаёт задачу. POST /api/v1/tasks. Body: CreateTaskRequest. Пользователь из JWT — создатель, должен быть в команде.
func (api *API) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := metadata.UserID(r.Context())
	if err != nil {
		mapError(w, r, err)
		return
	}

	var req dto.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mapError(w, r, err)
		return
	}

	if err := validation.Validator.Struct(req); err != nil {
		mapError(w, r, err)
		return
	}

	teamID, input, err := taskconverter.CreateTaskRequestToDomain(req)
	if err != nil {
		mapError(w, r, err)
		return
	}

	task, err := api.taskService.Create(r.Context(), userID, teamID, input)
	if err != nil {
		mapError(w, r, err)
		return
	}

	pkghttp.WriteJSON(r.Context(), w, http.StatusCreated, taskconverter.TaskToResponse(task))
}
