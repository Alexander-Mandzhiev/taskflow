package reader

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/converter"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/resources"
)

// GetByID возвращает команду по id (без удалённых).
// При tx != nil запрос выполняется в транзакции.
func (r *repository) GetByID(ctx context.Context, tx *sqlx.Tx, teamID uuid.UUID) (model.Team, error) {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Question).
		Select("id", "name", "created_by", "created_at", "updated_at", "deleted_at").
		From("teams").
		Where(sq.Eq{"id": teamID.String()}).
		Where(sq.Expr("deleted_at IS NULL")).
		Limit(1).
		ToSql()
	if err != nil {
		return model.Team{}, fmt.Errorf("build get by id query: %w", err)
	}

	var row resources.TeamRow
	if tx != nil {
		err = tx.GetContext(ctx, &row, query, args...)
	} else {
		err = r.readPool.GetContext(ctx, &row, query, args...)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Team{}, model.ErrTeamNotFound
		}
		return model.Team{}, toDomainError(err)
	}

	team, err := converter.ToDomainTeam(row)
	if err != nil {
		return model.Team{}, err
	}
	return team, nil
}
