package reader

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/repository/converter"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/repository/resources"
)

// ListByUserID возвращает команды, где пользователь участник, с его ролью в каждой.
// Один запрос: teams JOIN team_members ON team_id WHERE team_members.user_id = userID.
func (r *repository) ListByUserID(ctx context.Context, tx *sqlx.Tx, userID string) ([]*model.TeamWithRole, error) {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Question).
		Select("t.id", "t.name", "t.created_by", "t.created_at", "t.updated_at", "t.deleted_at", "tm.role").
		From("teams t").
		InnerJoin("team_members tm ON tm.team_id = t.id").
		Where(sq.Eq{"tm.user_id": userID}).
		Where(sq.Expr("t.deleted_at IS NULL")).
		OrderBy("t.name").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build list by user id query: %w", err)
	}

	var rows []resources.TeamWithRoleRow
	if tx != nil {
		err = tx.SelectContext(ctx, &rows, query, args...)
	} else {
		err = r.readPool.SelectContext(ctx, &rows, query, args...)
	}
	if err != nil {
		return nil, toDomainError(err)
	}

	out := make([]*model.TeamWithRole, 0, len(rows))
	for i := range rows {
		item, err := converter.ToDomainTeamWithRole(rows[i])
		if err != nil {
			return nil, fmt.Errorf("convert row %d: %w", i, err)
		}
		out = append(out, &item)
	}
	return out, nil
}
