package dto

// TeamTaskStatsResponse — статистика по команде (отчёт: участники, задачи done за период).
type TeamTaskStatsResponse struct {
	TeamID         string `json:"team_id"`
	TeamName       string `json:"team_name"`
	MemberCount    int    `json:"member_count"`
	DoneTasksCount int    `json:"done_tasks_count"`
}

// TeamTaskStatsListResponse — список статистики по командам.
type TeamTaskStatsListResponse struct {
	Items []TeamTaskStatsResponse `json:"items"`
}
