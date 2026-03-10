package writer

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

// Create создаёт запись приглашения. Поля inv (id, token, status, expires_at и т.д.) заданы вызывающим.
func (r *repository) Create(ctx context.Context, tx *sqlx.Tx, inv model.TeamInvitation) error {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Question).
		Insert("team_invitations").
		Columns("id", "team_id", "email", "role", "invited_by", "status", "token", "expires_at").
		Values(
			inv.ID.String(),
			inv.TeamID.String(),
			inv.Email,
			inv.Role,
			inv.InvitedBy.String(),
			inv.Status,
			inv.Token,
			inv.ExpiresAt,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("build create invitation query: %w", err)
	}

	exec := r.writePool.ExecContext
	if tx != nil {
		exec = tx.ExecContext
	}
	_, err = exec(ctx, query, args...)
	if err != nil {
		return toDomainError(err)
	}
	return nil
}
