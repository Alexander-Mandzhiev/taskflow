package list

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

const keyPrefix = "tasks:list"

// Key возвращает ключ Redis для закешированной страницы списка задач.
// Один ключ = одна страница (team_id + фильтр: status, assignee_id, limit, offset).
func Key(teamID uuid.UUID, filter model.TaskListFilter) string {
	return fmt.Sprintf("%s:%s:%s", keyPrefix, teamID.String(), filterSuffix(filter))
}

// PrefixForTeam возвращает префикс для инвалидации всех страниц списка команды (DelByPrefix).
func PrefixForTeam(teamID uuid.UUID) string {
	return fmt.Sprintf("%s:%s:", keyPrefix, teamID.String())
}

func filterSuffix(filter model.TaskListFilter) string {
	status := ""
	if filter.Status != nil {
		status = *filter.Status
	}
	assignee := ""
	if filter.AssigneeID != nil {
		assignee = filter.AssigneeID.String()
	}
	return fmt.Sprintf("s_%s_a_%s_l_%d_o_%d", status, assignee, filter.Limit, filter.Offset)
}
