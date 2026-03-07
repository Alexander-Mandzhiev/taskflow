package resources

import "time"

// TaskRow — строка таблицы tasks для чтения (поля с db-тегами для sqlx).
type TaskRow struct {
	ID          string     `db:"id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Status      string     `db:"status"`
	AssigneeID  *string    `db:"assignee_id"`
	TeamID      string     `db:"team_id"`
	CreatedBy   string     `db:"created_by"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	CompletedAt *time.Time `db:"completed_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}
