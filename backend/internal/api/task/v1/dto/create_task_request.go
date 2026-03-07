package dto

// CreateTaskRequest — запрос на создание задачи.
type CreateTaskRequest struct {
	TeamID      string  `json:"team_id" validate:"required,uuid"`
	Title       string  `json:"title" validate:"required,max=500"`
	Description string  `json:"description"`
	Status      string  `json:"status" validate:"omitempty,oneof=todo in_progress done"`
	AssigneeID  *string `json:"assignee_id" validate:"omitempty,uuid"`
}
