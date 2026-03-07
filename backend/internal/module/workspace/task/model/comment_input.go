package model

// CommentInput — входные данные для создания комментария к задаче (тело запроса).
// Используется при создании комментария; taskID и userID передаются в сигнатуре сервиса.
type CommentInput struct {
	Content string
}
