package resources

import "time"

// TeamInvitationRow — строка таблицы team_invitations для чтения (поля с db-тегами для sqlx).
type TeamInvitationRow struct {
	ID        string    `db:"id"`
	TeamID    string    `db:"team_id"`
	Email     string    `db:"email"`
	Role      string    `db:"role"`
	InvitedBy string    `db:"invited_by"`
	Status    string    `db:"status"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
