package writer

import (
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database"
)

// toDomainError преобразует ошибку БД в доменную ошибку модуля task.
func toDomainError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, model.ErrTaskNotFound) ||
		errors.Is(err, model.ErrInvalidStatus) ||
		errors.Is(err, model.ErrAssigneeNotInTeam) ||
		errors.Is(err, model.ErrTemporaryFailure) ||
		errors.Is(err, model.ErrInternal) {
		return err
	}

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		switch int(mysqlErr.Number) {
		case database.MySQLDuplicateEntry:
			return model.ErrInternal
		case database.MySQLForeignKeyConstraint:
			return model.ErrTaskNotFound
		case database.MySQLBadNull:
			return model.ErrInternal
		case database.MySQLDeadlock, database.MySQLLockWaitTimeout:
			return model.ErrTemporaryFailure
		}
	}

	return fmt.Errorf("%w: %w", model.ErrInternal, err)
}
