package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
)

// InvitationReaderRepository — чтение из таблицы team_invitations.
type InvitationReaderRepository interface {
	// GetPendingByTeamAndEmail возвращает приглашение со статусом pending для пары (team_id, email), если есть.
	// Если нет — (nil, model.ErrInvitationNotFound).
	GetPendingByTeamAndEmail(ctx context.Context, tx *sqlx.Tx, teamID, email string) (*model.TeamInvitation, error)
}

// InvitationWriterRepository — запись в таблицу team_invitations.
type InvitationWriterRepository interface {
	// Create создаёт запись приглашения (id, token, status=pending задаются вызывающим).
	Create(ctx context.Context, tx *sqlx.Tx, inv *model.TeamInvitation) error
}
