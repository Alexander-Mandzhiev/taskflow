package model

import "github.com/google/uuid"

// TeamTaskStats — результат сложного запроса (а): для каждой команды — название, кол-во участников, кол-во задач done за период.
type TeamTaskStats struct {
	TeamID         uuid.UUID
	TeamName       string
	MemberCount    int
	DoneTasksCount int // задачи в статусе done за последние 7 дней (или переданный период)
}

// TeamTopCreator — результат сложного запроса (б): топ-N пользователей по количеству созданных задач в команде за период.
type TeamTopCreator struct {
	TeamID       uuid.UUID
	UserID       uuid.UUID
	Rank         int // 1, 2, 3 (оконная функция)
	CreatedCount int64
}
