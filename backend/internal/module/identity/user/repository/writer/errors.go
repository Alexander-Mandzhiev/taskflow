package writer

import (
	"errors"

	"github.com/go-sql-driver/mysql"

	"mkk/internal/module/identity/user/repository/resources"
)

func isDuplicateKeyError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == resources.MySQLDuplicateEntry
	}
	return false
}
