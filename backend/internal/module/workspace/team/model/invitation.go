package model

import (
	"time"

	"github.com/google/uuid"
)

// Статусы приглашения (team_invitations.status).
const (
	InvitationStatusPending  = "pending"
	InvitationStatusAccepted = "accepted"
	InvitationStatusDeclined = "declined"
	InvitationStatusExpired  = "expired"
)

// TeamInvitation — приглашение в команду (таблица team_invitations).
// Добавление в team_members происходит только при принятии (AcceptInvitation).
type TeamInvitation struct {
	ID        uuid.UUID
	TeamID    uuid.UUID
	Email     string
	Role      string
	InvitedBy uuid.UUID
	Status    string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
