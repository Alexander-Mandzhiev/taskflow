package resources

// TeamWithRoleRow — результат запроса teams JOIN team_members: команда + роль пользователя (для ListByUserID).
// Колонки: id, name, created_by, created_at, updated_at, deleted_at — из teams; role — из team_members.
type TeamWithRoleRow struct {
	TeamRow
	Role string `db:"role"`
}
