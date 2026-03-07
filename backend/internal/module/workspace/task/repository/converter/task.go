package converter

import (
	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/resources"
)

// ToDomainTask преобразует строку БД (TaskRow) в доменную модель Task.
func ToDomainTask(r resources.TaskRow) (model.Task, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return model.Task{}, err
	}
	teamID, err := uuid.Parse(r.TeamID)
	if err != nil {
		return model.Task{}, err
	}
	createdBy, err := uuid.Parse(r.CreatedBy)
	if err != nil {
		return model.Task{}, err
	}
	var assigneeID *uuid.UUID
	if r.AssigneeID != nil && *r.AssigneeID != "" {
		a, err := uuid.Parse(*r.AssigneeID)
		if err != nil {
			return model.Task{}, err
		}
		assigneeID = &a
	}
	return model.Task{
		ID:          id,
		Title:       r.Title,
		Description: r.Description,
		Status:      r.Status,
		AssigneeID:  assigneeID,
		TeamID:      teamID,
		CreatedBy:   createdBy,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		CompletedAt: r.CompletedAt,
		DeletedAt:   r.DeletedAt,
	}, nil
}

// ToRepoTaskCreateInput преобразует доменный TaskInput в ресурс репозитория для INSERT. teamID передаётся в сигнатуре.
// При недопустимом переданном статусе возвращает model.ErrInvalidStatus.
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
