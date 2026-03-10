package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

// TeamService — слой сервиса команд и участников.
// Транзакции открываются внутри сервиса (txmanager.WithTx); вызывающий не передаёт tx.
// ID передаются типобезопасно как uuid.UUID. Create: при нулевом input возвращает model.ErrNilInput.
type TeamService interface {
	Create(ctx context.Context, input model.TeamInput, ownerUserID uuid.UUID) (model.Team, error)
	// GetByID возвращает команду с участниками только если userID — участник команды; иначе ErrForbidden.
	GetByID(ctx context.Context, teamID, userID uuid.UUID) (model.TeamWithMembers, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.TeamWithRole, error)
	GetMember(ctx context.Context, teamID, userID uuid.UUID) (model.TeamMember, error)
	// InviteByEmail создаёт приглашение (запись в team_invitations). Проверяет права (owner/admin), что пользователь не в команде и нет pending-приглашения. Отправка письма — отдельный сервис (позже).
	InviteByEmail(ctx context.Context, teamID, inviterUserID uuid.UUID, inviteeEmail, role string) (model.TeamInvitation, error)
	// AcceptInvitation принимает приглашение по токену из ссылки: проверяет токен, срок, что пользователь — приглашённый, добавляет в команду с ролью из приглашения. Пока не реализовано (ErrNotImplemented).
	AcceptInvitation(ctx context.Context, token string, userID uuid.UUID) (model.TeamMember, error)
}
