package task_v1

import (
	"net/http"
	"net/url"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/task/v1/dto"
	taskconverter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/converter"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/util"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

// List возвращает список задач с фильтрацией и пагинацией. GET /api/v1/tasks?team_id=...&status=...&assignee_id=...&limit=...&offset=...
func (api *API) List(w http.ResponseWriter, r *http.Request) {
	userID, err := metadata.UserID(r.Context())
	if err != nil {
		mapError(w, r, err)
		return
	}

	filter, err := parseTaskListFilter(r.URL.Query())
	if err != nil {
		mapError(w, r, err)
		return
	}

	items, total, err := api.taskService.List(r.Context(), userID, filter)
	if err != nil {
		mapError(w, r, err)
		return
	}

	pkghttp.WriteJSON(r.Context(), w, http.StatusOK, dto.TaskListResponse{
		Items:  taskconverter.TasksToResponse(items),
		Total:  total,
		Limit:  filter.Limit,
		Offset: filter.Offset,
	})
}

// parseTaskListFilter собирает TaskListFilter из query-параметров. Использует pkg/util для парсинга.
func parseTaskListFilter(q url.Values) (model.TaskListFilter, error) {
	teamIDStr := q.Get("team_id")
	if teamIDStr == "" {
		return model.TaskListFilter{}, model.ErrTeamIDRequired
	}
	teamID, err := util.ParseUUID(teamIDStr)
	if err != nil {
		return model.TaskListFilter{}, model.ErrForbidden
	}

	limit := defaultLimit
	if l := q.Get("limit"); l != "" {
		limit, err = util.ParseInt(l)
		if err != nil || limit <= 0 {
			return model.TaskListFilter{}, model.ErrPaginationRequired
		}
		if limit > maxLimit {
			limit = maxLimit
		}
	}

	offset := 0
	if o := q.Get("offset"); o != "" {
		n, err := util.ParseInt(o)
		if err == nil && n >= 0 {
			offset = n
		}
	}

	filter := model.TaskListFilter{
		TeamID: &teamID,
		Limit:  limit,
		Offset: offset,
	}
	if s := q.Get("status"); s != "" {
		filter.Status = &s
	}
	if a := q.Get("assignee_id"); a != "" {
		aid, err := util.ParseUUID(a)
		if err != nil {
			return model.TaskListFilter{}, model.ErrInvalidAssigneeID
		}
		filter.AssigneeID = &aid
	}
	return filter, nil
}
