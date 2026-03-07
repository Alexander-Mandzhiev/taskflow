package writer

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
)

// AddMember добавляет участника в команду. При нарушении uk возвращает model.ErrAlreadyMember.
func (r *repository) AddMember(ctx context.Context, tx *sqlx.Tx, teamID, userID, role string) (*model.TeamMember, error) {
	id := uuid.New().String()

	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Question).
		Insert("team_members").
		Columns("id", "user_id", "team_id", "role").
		Values(id, userID, teamID, role).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build add member query: %w", err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, toDomainError(err)
	}

	return r.selectByID(ctx, tx, id)
}
