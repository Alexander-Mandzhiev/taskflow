package task_v1

import (
	"net/http"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/task/v1/dto"
	taskconverter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/converter"
	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
)

// ReportInvalidAssignees возвращает задачи, где assignee не является членом команды задачи (валидация целостности).
// GET /api/v1/reports/invalid-assignees.
func (api *API) ReportInvalidAssignees(w http.ResponseWriter, r *http.Request) {
	userID, err := metadata.UserID(r.Context())
	if err != nil {
		mapError(w, r, err)
		return
	}

	tasks, err := api.reportService.TasksWithInvalidAssignee(r.Context(), userID)
	if err != nil {
		mapError(w, r, err)
		return
	}

	pkghttp.WriteJSON(r.Context(), w, http.StatusOK, dto.InvalidAssigneesListResponse{
		Items: taskconverter.TasksToResponse(tasks),
	})
}
