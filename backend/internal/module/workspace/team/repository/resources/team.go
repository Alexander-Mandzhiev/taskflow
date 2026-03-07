package resources

import "time"

// TeamRow — строка таблицы teams для чтения (поля с db-тегами для sqlx).
type TeamRow struct {
	ID        string     `db:"id"`
	Name      string     `db:"name"`
	CreatedBy string     `db:"created_by"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
