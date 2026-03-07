package reader

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

// GetPendingByTeamAndEmail возвращает приглашение со статусом pending для (team_id, email) или ErrInvitationNotFound.
func (r *repository) GetPendingByTeamAndEmail(ctx context.Context, tx *sqlx.Tx, teamID, email string) (*model2.TeamInvitation, error) {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Question).
		Select("id", "team_id", "email", "role", "invited_by", "status", "token", "expires_at", "created_at", "updated_at").
		From("team_invitations").
		Where(sq.Eq{"team_id": teamID, "email": email, "status": model2.InvitationStatusPending}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get pending invitation query: %w", err)
	}

	var row resources.TeamInvitationRow
	if tx != nil {
		err = tx.GetContext(ctx, &row, query, args...)
	} else {
		err = r.readPool.GetContext(ctx, &row, query, args...)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model2.ErrInvitationNotFound
		}
		return nil, toDomainError(err)
	}

	inv, err := converter.ToDomainTeamInvitation(row)
	if err != nil {
		return nil, err
	}
	return &inv, nil
}
