package writer

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// Update обновляет изменяемые поля задачи (title, description, status, assignee_id). При отсутствии — model.ErrTaskNotFound.
func (r *repository) Update(ctx context.Context, tx *sqlx.Tx, taskID uuid.UUID, input model.TaskInput) error {
	if tx == nil {
		return model.ErrTxRequired
	}

	var assigneeID interface{}
	if input.AssigneeID != nil {
		assigneeID = input.AssigneeID.String()
	}
	var completedAt interface{}
	if input.Status == model.TaskStatusDone {
		completedAt = sq.Expr("NOW()")
	} else {
		completedAt = nil
	}

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Question).
		Update("tasks").
		Set("title", input.Title).
		Set("description", input.Description).
		Set("status", input.Status).
		Set("assignee_id", assigneeID).
		Set("updated_at", sq.Expr("NOW()")).
		Set("completed_at", completedAt).
		Where(sq.Eq{"id": taskID.String()}).
		Where(sq.Expr("deleted_at IS NULL"))

	query, args, err := builder.ToSql()
	if err != nil {
		return toDomainError(err)
	}

	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return toDomainError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return model.ErrTaskNotFound
	}
	return nil
}
