package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

// TeamReaderRepository — чтение из таблицы teams (и JOIN с team_members для списка по user).
// tx != nil — запрос в транзакции; tx == nil — вне транзакции.
type TeamReaderRepository interface {
	// GetByID возвращает команду по id. При отсутствии — (nil, model.ErrTeamNotFound).
	GetByID(ctx context.Context, tx *sqlx.Tx, teamID uuid.UUID) (*model.Team, error)

	// ListByUserID возвращает команды, где пользователь является участником, с его ролью в каждой.
	// Один запрос: teams JOIN team_members ON team_id WHERE team_members.user_id = userID.
	ListByUserID(ctx context.Context, tx *sqlx.Tx, userID uuid.UUID) ([]*model.TeamWithRole, error)
}

// TeamWriterRepository — запись в таблицу teams.
// Мутации выполняются в транзакции (tx из txmanager.WithTx).
type TeamWriterRepository interface {
	// Create создаёт запись в teams. createdBy — user_id создателя (добавление в team_members делает адаптер).
	Create(ctx context.Context, tx *sqlx.Tx, input *model.TeamInput, createdBy uuid.UUID) (*model.Team, error)
}
