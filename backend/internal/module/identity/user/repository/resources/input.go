// Package resources — модели репозитория для таблицы users (без полей, проставляемых БД).
package resources

// UserInput — профильные данные для INSERT/UPDATE в users.
// password_hash передаётся отдельно (Create — аргумент, UpdatePasswordHash — отдельный метод).
type UserInput struct {
	Email string
	Name  string
}
