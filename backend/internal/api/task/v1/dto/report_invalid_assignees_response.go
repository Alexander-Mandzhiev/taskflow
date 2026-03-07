package dto

// InvalidAssigneesListResponse — задачи, где assignee не в команде (валидация целостности).
type InvalidAssigneesListResponse struct {
	Items []TaskResponse `json:"items"`
}
