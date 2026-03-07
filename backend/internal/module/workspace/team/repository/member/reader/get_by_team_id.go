package reader

import (
	"context"
	"fmt"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/converter"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/resources"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

// GetByTeamID возвращает всех участников команды по team_id.
func (r *repository) GetByTeamID(ctx context.Context, tx *sqlx.Tx, teamID string) ([]*model.TeamMember, error) {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Question).
		Select("id", "user_id", "team_id", "role", "created_at").
		From("team_members").
		Where(sq.Eq{"team_id": teamID}).
		OrderBy("created_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get by team id query: %w", err)
	}

	var rows []resources.TeamMemberRow
	if tx != nil {
		err = tx.SelectContext(ctx, &rows, query, args...)
	} else {
		err = r.readPool.SelectContext(ctx, &rows, query, args...)
	}
	if err != nil {
		return nil, toDomainError(err)
	}

	out := make([]*model.TeamMember, 0, len(rows))
	for i := range rows {
		member, err := converter.ToDomainTeamMember(rows[i])
		if err != nil {
			return nil, fmt.Errorf("convert row %d: %w", i, err)
		}
		out = append(out, &member)
	}
	return out, nil
}
