package model

// UserInput — профильные данные пользователя (создание и обновление профиля).
// PasswordHash сюда не входит — при создании передаётся отдельным аргументом,
// смена пароля выполняется через отдельный метод UpdatePasswordHash.
type UserInput struct {
	Email string
	Name  string
}
