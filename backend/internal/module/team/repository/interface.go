package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
)

// TeamRepository — единая точка доступа к данным команд и участников (контракт для адаптера).
// tx: при tx != nil все операции в транзакции; при tx == nil — вне транзакции.
type TeamRepository interface {
	// Create создаёт только запись в teams (created_by = ownerUserID). Добавление owner в team_members — зона ответственности сервиса (AddMember).
	Create(ctx context.Context, tx *sqlx.Tx, input *model.TeamInput, ownerUserID string) (*model.Team, error)

	// GetByID возвращает команду с участниками. При отсутствии — (nil, model.ErrTeamNotFound). Оба чтения в одном tx при tx != nil.
	GetByID(ctx context.Context, tx *sqlx.Tx, teamID string) (*model.TeamWithMembers, error)

	// ListByUserID — список команд, где пользователь участник, с его ролью в каждой (GET /api/v1/teams). Без списка членов.
	ListByUserID(ctx context.Context, tx *sqlx.Tx, userID string) ([]*model.TeamWithRole, error)

	// GetMember — участник по (team_id, user_id). Для проверки прав (owner/admin) и «уже в команде». При отсутствии — (nil, model.ErrMemberNotFound).
	GetMember(ctx context.Context, tx *sqlx.Tx, teamID, userID string) (*model.TeamMember, error)

	// AddMember добавляет пользователя в команду с указанной ролью (для invite). При дубликате — ошибка (model.ErrAlreadyMember или от БД).
	AddMember(ctx context.Context, tx *sqlx.Tx, teamID, userID, role string) (*model.TeamMember, error)
}
