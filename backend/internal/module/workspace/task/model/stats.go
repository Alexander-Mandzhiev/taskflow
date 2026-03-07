package model

import "github.com/google/uuid"

// TeamTaskStats — результат отчёта: для каждой команды — название, кол-во участников, кол-во задач done за период.
type TeamTaskStats struct {
	TeamID         uuid.UUID
	TeamName       string
	MemberCount    int
	DoneTasksCount int
}

// TeamTopCreator — результат отчёта: топ-N пользователей по количеству созданных задач в команде за период.
type TeamTopCreator struct {
	TeamID       uuid.UUID
	UserID       uuid.UUID
	Rank         int
	CreatedCount int64
}
