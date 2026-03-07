package reader

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/converter"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/resources"
)

// List возвращает список задач с фильтром и пагинацией. total — количество без LIMIT.
func (r *repository) List(ctx context.Context, tx *sqlx.Tx, filter *model.TaskListFilter, pagination *model.TaskPagination) ([]*model.Task, int, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Question)
	where := sq.Expr("deleted_at IS NULL")
	if filter != nil {
		if filter.TeamID != nil {
			where = sq.And{where, sq.Eq{"team_id": filter.TeamID.String()}}
		}
		if filter.Status != nil && *filter.Status != "" {
			where = sq.And{where, sq.Eq{"status": *filter.Status}}
		}
		if filter.AssigneeID != nil {
			where = sq.And{where, sq.Eq{"assignee_id": filter.AssigneeID.String()}}
		}
	}

	// COUNT для total
	countQuery, countArgs, err := builder.Select("COUNT(*)").
		From("tasks").
		Where(where).
		ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build count query: %w", err)
	}

	var total int
	if tx != nil {
		err = tx.GetContext(ctx, &total, countQuery, countArgs...)
	} else {
		err = r.readPool.GetContext(ctx, &total, countQuery, countArgs...)
	}
	if err != nil {
		return nil, 0, toDomainError(err)
	}

	// SELECT с пагинацией
	limit, offset := 20, 0
	if pagination != nil {
		if pagination.Limit > 0 {
			limit = pagination.Limit
		}
		if pagination.Offset > 0 {
			offset = pagination.Offset
		}
	}

	listQuery, listArgs, err := builder.
		Select("id", "title", "description", "status", "assignee_id", "team_id", "created_by", "created_at", "updated_at", "deleted_at").
		From("tasks").
		Where(where).
		OrderBy("updated_at DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build list query: %w", err)
	}

	var rows []resources.TaskRow
	if tx != nil {
		err = tx.SelectContext(ctx, &rows, listQuery, listArgs...)
	} else {
		err = r.readPool.SelectContext(ctx, &rows, listQuery, listArgs...)
	}
	if err != nil {
		return nil, 0, toDomainError(err)
	}

	out := make([]*model.Task, 0, len(rows))
	for i := range rows {
		task, err := converter.ToDomainTask(rows[i])
		if err != nil {
			return nil, 0, fmt.Errorf("convert row %d: %w", i, err)
		}
		out = append(out, &task)
	}
	return out, total, nil
}
