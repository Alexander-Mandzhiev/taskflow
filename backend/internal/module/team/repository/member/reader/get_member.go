package reader

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

// GetMember возвращает участника по паре (team_id, user_id). При отсутствии — model.ErrMemberNotFound.
func (r *repository) GetMember(ctx context.Context, tx *sqlx.Tx, teamID, userID string) (*model.TeamMember, error) {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Question).
		Select("id", "user_id", "team_id", "role", "created_at").
		From("team_members").
		Where(sq.Eq{"team_id": teamID, "user_id": userID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get member query: %w", err)
	}

	var row resources.TeamMemberRow
	if tx != nil {
		err = tx.GetContext(ctx, &row, query, args...)
	} else {
		err = r.readPool.GetContext(ctx, &row, query, args...)
	}
	if err != nil {
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
