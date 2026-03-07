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

// GetPendingByTeamAndEmail возвращает приглашение со статусом pending для (team_id, email) или ErrInvitationNotFound.
func (r *repository) GetPendingByTeamAndEmail(ctx context.Context, tx *sqlx.Tx, teamID, email string) (*model.TeamInvitation, error) {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Question).
		Select("id", "team_id", "email", "role", "invited_by", "status", "token", "expires_at", "created_at", "updated_at").
		From("team_invitations").
		Where(sq.Eq{"team_id": teamID, "email": email, "status": model.InvitationStatusPending}).
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
			return nil, model.ErrInvitationNotFound
		}
		return nil, toDomainError(err)
	}

	inv, err := converter.ToDomainTeamInvitation(row)
	if err != nil {
		return nil, err
	}
	return &inv, nil
}
