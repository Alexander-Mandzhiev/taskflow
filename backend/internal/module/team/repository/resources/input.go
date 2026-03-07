package resources

// TeamInput — данные для INSERT в таблицу teams (создание команды).
// created_by передаётся отдельным аргументом в Create.
type TeamInput struct {
	Name string
}
