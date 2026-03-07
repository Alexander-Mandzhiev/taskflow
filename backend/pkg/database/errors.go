package database

// Коды MySQL для маппинга в доменные ошибки модулей (user, team и т.д.).
const (
	MySQLDuplicateEntry       = 1062 // ER_DUP_ENTRY
	MySQLForeignKeyConstraint = 1452 // ER_NO_REFERENCED_ROW_2
	MySQLBadNull              = 1048 // ER_BAD_NULL_ERROR
	MySQLDeadlock             = 1213 // ER_LOCK_DEADLOCK
	MySQLLockWaitTimeout      = 1205 // ER_LOCK_WAIT_TIMEOUT
)
