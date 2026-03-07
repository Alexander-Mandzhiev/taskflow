package resources

import "time"

// TaskHistoryRow — строка таблицы task_history для чтения/записи (db-теги для sqlx).
// ChangedAt при INSERT проставляется в БД (DEFAULT now()) или вызывающим при формировании model.TaskHistory.
type TaskHistoryRow struct {
	ID        string    `db:"id"`
	TaskID    string    `db:"task_id"`
	ChangedBy string    `db:"changed_by"`
	FieldName string    `db:"field_name"`
	OldValue  string    `db:"old_value"`
	NewValue  string    `db:"new_value"`
	ChangedAt time.Time `db:"changed_at"`
}
