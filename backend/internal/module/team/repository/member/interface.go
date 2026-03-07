package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
)

// MemberReaderRepository — чтение из таблицы team_members.
type MemberReaderRepository interface {
	// GetByTeamID возвращает всех участников команды по team_id.
	GetByTeamID(ctx context.Context, tx *sqlx.Tx, teamID string) ([]*model.TeamMember, error)

	// GetMember возвращает запись участника по паре (team_id, user_id).
	// Для проверки прав (owner/admin) и проверки «уже в команде». При отсутствии — (nil, model.ErrMemberNotFound).
	GetMember(ctx context.Context, tx *sqlx.Tx, teamID, userID string) (*model.TeamMember, error)
}

// MemberWriterRepository — запись в таблицу team_members.
type MemberWriterRepository interface {
	// AddMember добавляет участника в команду (user_id, team_id, role).
	// При нарушении uk_team_members_user_team репозиторий возвращает ошибку (сервис может маппить в model.ErrAlreadyMember).
	AddMember(ctx context.Context, tx *sqlx.Tx, teamID, userID, role string) (*model.TeamMember, error)
}
