package dto

// CommentResponse — комментарий к задаче в ответе API.
type CommentResponse struct {
	ID        string `json:"id"`
	TaskID    string `json:"task_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CommentListResponse — список комментариев к задаче.
type CommentListResponse struct {
	Items []CommentResponse `json:"items"`
}
