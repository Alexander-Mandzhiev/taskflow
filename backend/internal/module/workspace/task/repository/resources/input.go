package resources

// TaskCreateInput — данные для INSERT в таблицу tasks.
// id, created_by передаются отдельно в Create; created_at/updated_at — в БД.
type TaskCreateInput struct {
	Title       string
	Description string
	Status      string
	AssigneeID  *string // NULL если не назначен
	TeamID      string
}
