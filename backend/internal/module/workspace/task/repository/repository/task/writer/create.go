package writer

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/converter"
)

// Create создаёт запись в tasks. teamID и createdBy — в сигнатуре. Мутация только в транзакции (tx != nil).
func (r *repository) Create(ctx context.Context, tx *sqlx.Tx, teamID uuid.UUID, input *model.TaskInput, createdBy uuid.UUID) (*model.Task, error) {
	if tx == nil {
		return nil, model.ErrTxRequired
	}
	in, err := converter.ToRepoTaskCreateInput(teamID, input)
	if err != nil {
		return nil, err
	}

	id := uuid.New().String()
	var assigneeID interface{}
	if in.AssigneeID != nil {
		assigneeID = *in.AssigneeID
	}
	var completedAt interface{}
	if in.Status == model.TaskStatusDone {
		completedAt = sq.Expr("NOW()")
	} else {
		completedAt = nil
	}

	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Question).
		Insert("tasks").
		Columns("id", "title", "description", "status", "assignee_id", "team_id", "created_by", "created_at", "updated_at", "completed_at").
		Values(id, in.Title, in.Description, in.Status, assigneeID, in.TeamID, createdBy.String(), sq.Expr("NOW()"), sq.Expr("NOW()"), completedAt).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build create query: %w", err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, toDomainError(err)
	}

	return r.selectByID(ctx, tx, id)
}
