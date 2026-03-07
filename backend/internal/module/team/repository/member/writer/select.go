package writer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/repository/converter"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/repository/resources"
)

// selectByID читает участника по id внутри текущей транзакции (после AddMember).
func (r *repository) selectByID(ctx context.Context, tx *sqlx.Tx, id string) (*model.TeamMember, error) {
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
			return nil, model.ErrMemberNotFound
		}
		return nil, toDomainError(err)
	}

	member, err := converter.ToDomainTeamMember(row)
	if err != nil {
		return nil, err
	}
	return &member, nil
}
