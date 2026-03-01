package resources

import "time"

// UserRow — строка users для чтения (все поля из БД).
type UserRow struct {
	ID           string     `db:"id"`
	Email        string     `db:"email"`
	Name         string     `db:"name"`
	PasswordHash string     `db:"password_hash"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}
