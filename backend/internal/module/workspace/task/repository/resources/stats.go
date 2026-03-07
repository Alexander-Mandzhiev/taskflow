package resources

// TeamTaskStatsRow — строка результата запроса отчёта: по команде — название, кол-во участников, кол-во задач done за период.
type TeamTaskStatsRow struct {
	TeamID         string `db:"team_id"`
	TeamName       string `db:"team_name"`
	MemberCount    int    `db:"member_count"`
	DoneTasksCount int    `db:"done_tasks_count"`
}

// TeamTopCreatorRow — строка результата отчёта: топ-N по созданным задачам в команде за период (оконная функция).
type TeamTopCreatorRow struct {
	TeamID       string `db:"team_id"`
	UserID       string `db:"user_id"`
	Rank         int    `db:"rank"`
	CreatedCount int64  `db:"created_count"`
}
