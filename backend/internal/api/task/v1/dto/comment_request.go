package dto

// CreateCommentRequest — запрос на создание комментария к задаче.
type CreateCommentRequest struct {
	Content string `json:"content" validate:"required,max=10000"`
}
