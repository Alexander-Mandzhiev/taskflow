package writer

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository/converter"
)

// Create создаёт пользователя и возвращает сохранённую сущность.
// UUID генерируется на стороне Go; created_at/updated_at проставляются MySQL (DEFAULT CURRENT_TIMESTAMP).
// Вызывается только внутри txmanager.WithTx, tx передаётся явно.
func (r *repository) Create(ctx context.Context, tx *sqlx.Tx, input model.UserInput, passwordHash string) (model.User, error) {
	in := converter.ToRepoInput(input)
	id := uuid.New().String()

	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Question).
		Insert("users").
		Columns("id", "email", "name", "password_hash").
		Values(id, in.Email, in.Name, passwordHash).
		ToSql()
	if err != nil {
		return model.User{}, fmt.Errorf("build create query: %w", err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return model.User{}, toDomainError(err)
	}

	return r.selectByID(ctx, tx, id)
}
