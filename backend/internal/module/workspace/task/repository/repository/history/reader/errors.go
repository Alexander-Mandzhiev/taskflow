package reader

import (
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database"
)

func toDomainError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, model.ErrTaskNotFound) ||
		errors.Is(err, model.ErrTemporaryFailure) ||
		errors.Is(err, model.ErrInternal) {
		return err
	}
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		switch int(mysqlErr.Number) {
		case database.MySQLDeadlock, database.MySQLLockWaitTimeout:
			return model.ErrTemporaryFailure
		}
	}
	return fmt.Errorf("%w: %w", model.ErrInternal, err)
}
