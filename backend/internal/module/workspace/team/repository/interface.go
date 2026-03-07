package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	model2 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

// TeamRepository — единая точка доступа к данным команд и участников (контракт для адаптера).
// tx: при tx != nil все операции в транзакции; при tx == nil — вне транзакции. ID передаются как uuid.UUID.
type TeamRepository interface {
	// Create создаёт только запись в teams (created_by = ownerUserID). Добавление owner в team_members — зона ответственности сервиса (AddMember).
	Create(ctx context.Context, tx *sqlx.Tx, input *model2.TeamInput, ownerUserID uuid.UUID) (*model2.Team, error)

	// GetByID возвращает только команду. При отсутствии — (nil, model.ErrTeamNotFound).
	GetByID(ctx context.Context, tx *sqlx.Tx, teamID uuid.UUID) (*model2.Team, error)

	// GetMembersByTeamID возвращает участников команды по team_id. В одной tx с GetByID — консистентный снимок.
	GetMembersByTeamID(ctx context.Context, tx *sqlx.Tx, teamID uuid.UUID) ([]*model2.TeamMember, error)

	// ListByUserID — список команд, где пользователь участник, с его ролью в каждой (GET /api/v1/teams). Без списка членов.
	ListByUserID(ctx context.Context, tx *sqlx.Tx, userID uuid.UUID) ([]*model2.TeamWithRole, error)

	// GetMember — участник по (team_id, user_id). Для проверки прав (owner/admin) и «уже в команде». При отсутствии — (nil, model.ErrMemberNotFound).
	GetMember(ctx context.Context, tx *sqlx.Tx, teamID, userID uuid.UUID) (*model2.TeamMember, error)

	// AddMember добавляет пользователя в команду с указанной ролью (для invite). При дубликате — ошибка (model.ErrAlreadyMember или от БД).
	AddMember(ctx context.Context, tx *sqlx.Tx, teamID, userID uuid.UUID, role string) (*model2.TeamMember, error)

	// CreateInvitation создаёт запись приглашения в team_invitations (status=pending, token и expires_at заданы вызывающим).
	CreateInvitation(ctx context.Context, tx *sqlx.Tx, inv *model2.TeamInvitation) error

	// GetPendingInvitationByTeamAndEmail возвращает приглашение со статусом pending для (team_id, email) или model.ErrInvitationNotFound.
	GetPendingInvitationByTeamAndEmail(ctx context.Context, tx *sqlx.Tx, teamID uuid.UUID, email string) (*model2.TeamInvitation, error)
}
