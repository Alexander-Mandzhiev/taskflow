package dto

// UpdateTaskRequest — запрос на обновление задачи.
type UpdateTaskRequest struct {
	Title       string  `json:"title" validate:"required,max=500"`
	Description string  `json:"description"`
	Status      string  `json:"status" validate:"required,oneof=todo in_progress done"`
	AssigneeID  *string `json:"assignee_id" validate:"omitempty,uuid"`
}
