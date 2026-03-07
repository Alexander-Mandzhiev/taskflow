package converter

import (
	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/task/v1/dto"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// CreateTaskRequestToDomain конвертирует DTO запроса создания задачи в домен (TaskInput + teamID).
// teamID парсится из req.TeamID; assignee_id опционален.
func CreateTaskRequestToDomain(req dto.CreateTaskRequest) (teamID uuid.UUID, input *model.TaskInput, err error) {
	teamID, err = uuid.Parse(req.TeamID)
	if err != nil {
		return uuid.Nil, nil, err
	}
	input = &model.TaskInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	}
	if req.AssigneeID != nil && *req.AssigneeID != "" {
		aid, e := uuid.Parse(*req.AssigneeID)
		if e != nil {
			return uuid.Nil, nil, e
		}
		input.AssigneeID = &aid
	}
	return teamID, input, nil
}

// UpdateTaskRequestToDomain конвертирует DTO запроса обновления задачи в доменную модель.
func UpdateTaskRequestToDomain(req dto.UpdateTaskRequest) (*model.TaskInput, error) {
	input := &model.TaskInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	}
	if req.AssigneeID != nil && *req.AssigneeID != "" {
		aid, err := uuid.Parse(*req.AssigneeID)
		if err != nil {
			return nil, err
		}
		input.AssigneeID = &aid
	} else {
		input.AssigneeID = nil
	}
	return input, nil
}
