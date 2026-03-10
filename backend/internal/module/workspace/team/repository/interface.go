package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

// TeamRepository — доступ к данным команд (таблица teams). Контракт для адаптера.
// tx: при tx != nil — в транзакции; при tx == nil — вне.
type TeamRepository interface {
	Create(ctx context.Context, tx *sqlx.Tx, input model.TeamInput, ownerUserID uuid.UUID) (model.Team, error)
	GetByID(ctx context.Context, tx *sqlx.Tx, teamID uuid.UUID) (model.Team, error)
	ListByUserID(ctx context.Context, tx *sqlx.Tx, userID uuid.UUID) ([]model.TeamWithRole, error)
}

// MemberRepository — доступ к данным участников команд (таблица team_members). Контракт для адаптера.
type MemberRepository interface {
	GetMembersByTeamID(ctx context.Context, tx *sqlx.Tx, teamID uuid.UUID) ([]model.TeamMember, error)
	GetMember(ctx context.Context, tx *sqlx.Tx, teamID, userID uuid.UUID) (model.TeamMember, error)
	AddMember(ctx context.Context, tx *sqlx.Tx, teamID, userID uuid.UUID, role string) (model.TeamMember, error)
}

// InvitationRepository — доступ к приглашениям в команды (таблица team_invitations). Контракт для адаптера.
type InvitationRepository interface {
	CreateInvitation(ctx context.Context, tx *sqlx.Tx, inv model.TeamInvitation) error
	GetPendingInvitationByTeamAndEmail(ctx context.Context, tx *sqlx.Tx, teamID uuid.UUID, email string) (model.TeamInvitation, error)
}
