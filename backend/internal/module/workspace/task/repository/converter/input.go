package converter

import (
	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/resources"
)

// ToRepoTaskCreateInput преобразует доменный TaskInput в ресурс репозитория для INSERT. teamID передаётся в сигнатуре.
// При недопустимом переданном статусе (не пустая и не todo/in_progress/done) возвращает model.ErrInvalidStatus.
func ToRepoTaskCreateInput(teamID uuid.UUID, input *model.TaskInput) (resources.TaskCreateInput, error) {
	if input == nil {
		return resources.TaskCreateInput{}, nil
	}
	status := input.Status
	if status != "" && !model.IsValidTaskStatus(status) {
		return resources.TaskCreateInput{}, model.ErrInvalidStatus
	}
	if status == "" {
		status = model.TaskStatusTodo
	}
	out := resources.TaskCreateInput{
		TeamID:      teamID.String(),
		Title:       input.Title,
		Description: input.Description,
		Status:      status,
	}
	if input.AssigneeID != nil {
		s := input.AssigneeID.String()
		out.AssigneeID = &s
	}
	return out, nil
}

// ToRepoTaskHistory преобразует доменную запись TaskHistory в ресурс для INSERT в task_history.
func ToRepoTaskHistory(entry *model.TaskHistory) resources.TaskHistoryRow {
	if entry == nil {
		return resources.TaskHistoryRow{}
	}
	return resources.TaskHistoryRow{
		ID:        entry.ID.String(),
		TaskID:    entry.TaskID.String(),
		ChangedBy: entry.ChangedBy.String(),
		FieldName: entry.FieldName,
		OldValue:  entry.OldValue,
		NewValue:  entry.NewValue,
		ChangedAt: entry.ChangedAt,
	}
}
