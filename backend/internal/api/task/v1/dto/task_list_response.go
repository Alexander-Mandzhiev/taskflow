package dto

// TaskListResponse — список задач с пагинацией.
type TaskListResponse struct {
	Items  []TaskResponse `json:"items"`
	Total  int            `json:"total"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}
