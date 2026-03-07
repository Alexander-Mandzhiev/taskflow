package reader

import (
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database"
)

// toDomainError преобразует ошибку БД (сырой код MySQL) в доменную ошибку модуля user.
// Уже доменные ошибки возвращаются без изменений.
// Неизвестные коды MySQL и не-MySQL ошибки возвращаются как model.ErrInternal с сохранением цепочки.
func toDomainError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, model.ErrUserNotFound) ||
		errors.Is(err, model.ErrEmailDuplicate) ||
		errors.Is(err, model.ErrInvalidInput) ||
		errors.Is(err, model.ErrTemporaryFailure) ||
		errors.Is(err, model.ErrInternal) {
		return err
	}

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		switch int(mysqlErr.Number) {
		case database.MySQLDuplicateEntry:
			return model.ErrEmailDuplicate
		case database.MySQLForeignKeyConstraint:
			return model.ErrUserNotFound
		case database.MySQLBadNull:
			return model.ErrInvalidInput
		case database.MySQLDeadlock, database.MySQLLockWaitTimeout:
			return model.ErrTemporaryFailure
		}
	}

	return fmt.Errorf("%w: %w", model.ErrInternal, err)
}
