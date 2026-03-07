package resources

import "time"

// TaskCommentRow — строка таблицы task_comments для чтения (поля с db-тегами для sqlx).
type TaskCommentRow struct {
	ID        string     `db:"id"`
	TaskID    string     `db:"task_id"`
	UserID    string     `db:"user_id"`
	Content   string     `db:"content"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// TaskCommentCreateInput — данные для INSERT в таблицу task_comments.
// id генерируется в БД или вызывающим; created_at/updated_at — в БД.
type TaskCommentCreateInput struct {
	TaskID  string
	UserID  string
	Content string
}
