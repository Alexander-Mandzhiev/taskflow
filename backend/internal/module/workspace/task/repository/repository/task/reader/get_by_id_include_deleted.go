package reader

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/converter"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/resources"
)

// GetByIDIncludeDeleted возвращает задачу по id в том числе с deleted_at IS NOT NULL (для Restore).
func (r *repository) GetByIDIncludeDeleted(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID) (model.Task, error) {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Question).
		Select("id", "title", "description", "status", "assignee_id", "team_id", "created_by", "created_at", "updated_at", "completed_at", "deleted_at").
		From("tasks").
		Where(sq.Eq{"id": taskID.String()}).
		Limit(1).
		ToSql()
	if err != nil {
		return model.Task{}, fmt.Errorf("build get by id query: %w", err)
	}

	var row resources.TaskRow
	if tx != nil {
		err = tx.GetContext(ctx, &row, query, args...)
	} else {
		err = r.readPool.GetContext(ctx, &row, query, args...)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Task{}, model.ErrTaskNotFound
		}
		return model.Task{}, toDomainError(err)
	}

	task, err := converter.ToDomainTask(row)
	if err != nil {
		return model.Task{}, err
	}
	return task, nil
}
