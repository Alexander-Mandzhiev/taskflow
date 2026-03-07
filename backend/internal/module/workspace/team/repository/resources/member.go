package resources

import "time"

// TeamMemberRow — строка таблицы team_members для чтения (поля с db-тегами для sqlx).
type TeamMemberRow struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	TeamID    string    `db:"team_id"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
}
