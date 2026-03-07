package writer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	model2 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/converter"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/resources"
)

// selectByID читает команду по ID внутри текущей транзакции (после Create).
func (r *repository) selectByID(ctx context.Context, tx *sqlx.Tx, teamID string) (*model2.Team, error) {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Question).
		Select("id", "name", "created_by", "created_at", "updated_at", "deleted_at").
		From("teams").
		Where(sq.Eq{"id": teamID}).
		Where(sq.Expr("deleted_at IS NULL")).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select query: %w", err)
	}

	var row resources.TeamRow
	if err := tx.GetContext(ctx, &row, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model2.ErrTeamNotFound
		}
		return nil, toDomainError(err)
	}

	team, err := converter.ToDomainTeam(row)
	if err != nil {
		return nil, err
	}
	return &team, nil
}
