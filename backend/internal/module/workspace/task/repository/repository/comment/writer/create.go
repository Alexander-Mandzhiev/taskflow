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

// Create создаёт комментарий к задаче; id генерируется, created_at/updated_at — в БД.
func (r *repository) Create(ctx context.Context, tx *sqlx.Tx, taskID, userID uuid.UUID, content string) (*model.TaskComment, error) {
	id := uuid.New().String()
	in := converter.ToRepoTaskCommentCreateInput(taskID, userID, content)

	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Question).
		Insert("task_comments").
		Columns("id", "task_id", "user_id", "content", "created_at", "updated_at").
		Values(id, in.TaskID, in.UserID, in.Content, sq.Expr("NOW()"), sq.Expr("NOW()")).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build create query: %w", err)
	}

	if tx != nil {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = r.writePool.ExecContext(ctx, query, args...)
	}
	if err != nil {
		return nil, toDomainError(err)
	}

	return r.selectByID(ctx, tx, id)
}
