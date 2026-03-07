package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// ReportRepository — отчёты по задачам и командам (сложные запросы с JOIN/агрегацией).
// Не смешивается с CRUD задач: отдельный модуль, единый контракт для отчётов.
type ReportRepository interface {
	// TeamTaskStats — для каждой команды: название, кол-во участников, кол-во задач в статусе done за период since.
	TeamTaskStats(ctx context.Context, tx *sqlx.Tx, since time.Time) ([]*model.TeamTaskStats, error)

	// TopCreatorsByTeam — топ-N пользователей по количеству созданных задач в каждой команде за период.
	TopCreatorsByTeam(ctx context.Context, tx *sqlx.Tx, since time.Time, limit int) ([]*model.TeamTopCreator, error)

	// TasksWithInvalidAssignee — задачи, у которых assignee не является участником команды задачи.
	TasksWithInvalidAssignee(ctx context.Context, tx *sqlx.Tx) ([]*model.Task, error)
}
