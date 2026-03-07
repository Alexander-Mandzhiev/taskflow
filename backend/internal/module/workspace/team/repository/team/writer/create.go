package writer

import (
	"context"
	"fmt"

	model2 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/converter"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Create создаёт запись в teams. createdBy — user_id создателя.
// UUID генерируется в Go; created_at/updated_at проставляются MySQL.
func (r *repository) Create(ctx context.Context, tx *sqlx.Tx, input *model2.TeamInput, createdBy string) (*model2.Team, error) {
	in := converter.ToRepoTeamInput(input)
	id := uuid.New().String()

	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Question).
		Insert("teams").
		Columns("id", "name", "created_by").
		Values(id, in.Name, createdBy).
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
