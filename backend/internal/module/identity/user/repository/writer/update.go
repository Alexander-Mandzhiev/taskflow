package writer

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"mkk/internal/module/identity/user/model"
	"mkk/internal/module/identity/user/repository/converter"
)

// Update обновляет пользователя по ID (игнорирует удалённых) и возвращает обновлённую сущность.
// Вызывается только внутри txmanager.WithTx, tx передаётся явно.
// Валидация id выполняется в сервисном и API-слое.
func (r *repository) Update(ctx context.Context, tx *sqlx.Tx, id string, input *model.UserInput) (*model.User, error) {
	in := converter.ToRepoInput(input)

	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Question).
		Update("users").
		Set("email", in.Email).
		Set("name", in.Name).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id}).
		Where(sq.Expr("deleted_at IS NULL")).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build update query: %w", err)
	}

	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, model.ErrEmailDuplicate
		}
		return nil, fmt.Errorf("update exec: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return nil, model.ErrUserNotFound
	}

	return r.selectByID(ctx, tx, id)
}
