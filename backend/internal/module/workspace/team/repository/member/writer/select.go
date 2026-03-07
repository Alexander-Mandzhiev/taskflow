package writer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	model2 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/converter"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/resources"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

// selectByID читает участника по id внутри текущей транзакции (после AddMember).
func (r *repository) selectByID(ctx context.Context, tx *sqlx.Tx, id string) (*model2.TeamMember, error) {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Question).
		Select("id", "user_id", "team_id", "role", "created_at").
		From("team_members").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select query: %w", err)
	}

	var row resources.TeamMemberRow
	if err := tx.GetContext(ctx, &row, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model2.ErrMemberNotFound
		}
		return nil, toDomainError(err)
	}

	member, err := converter.ToDomainTeamMember(row)
	if err != nil {
		return nil, err
	}
	return &member, nil
}
