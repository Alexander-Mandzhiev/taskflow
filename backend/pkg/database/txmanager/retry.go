package txmanager

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

const defaultMaxRetries = 3

// MySQL error codes для serialization failures.
const (
	mysqlDeadlock        = 1213 // ER_LOCK_DEADLOCK
	mysqlLockWaitTimeout = 1205 // ER_LOCK_WAIT_TIMEOUT
)

// isSerializationError определяет, является ли ошибка serialization/deadlock failure,
// которую безопасно повторить.
func isSerializationError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == mysqlDeadlock ||
			mysqlErr.Number == mysqlLockWaitTimeout
	}
	return false
}
