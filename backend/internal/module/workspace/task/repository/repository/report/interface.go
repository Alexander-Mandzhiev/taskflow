package report

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// ReportReaderRepository — чтение отчётов по задачам и командам (сложные запросы).
type ReportReaderRepository interface {
	TeamTaskStats(ctx context.Context, tx *sqlx.Tx, since time.Time) ([]*model.TeamTaskStats, error)
	TopCreatorsByTeam(ctx context.Context, tx *sqlx.Tx, since time.Time, limit int) ([]*model.TeamTopCreator, error)
	TasksWithInvalidAssignee(ctx context.Context, tx *sqlx.Tx) ([]*model.Task, error)
}
