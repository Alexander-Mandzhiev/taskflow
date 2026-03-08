package member

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

// AddMember добавляет пользователя в команду с указанной ролью. При дубликате — ошибка (model.ErrAlreadyMember или от БД).
func (r *Adapter) AddMember(ctx context.Context, tx *sqlx.Tx, teamID, userID uuid.UUID, role string) (*model.TeamMember, error) {
	return r.memberWriter.AddMember(ctx, tx, teamID, userID, role)
}
