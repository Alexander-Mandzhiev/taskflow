package model

import (
	"time"

	"github.com/google/uuid"
)

// Роли участника команды (team_members.role).
const (
	RoleOwner  = "owner"
	RoleAdmin  = "admin"
	RoleMember = "member"
)

// TeamMember — участник команды (таблица team_members).
type TeamMember struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TeamID    uuid.UUID
	Role      string
	CreatedAt time.Time
}
