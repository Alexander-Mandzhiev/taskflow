package invitation

import (
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository"
	invitationRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/repository/invitation"
)

var _ repository.InvitationRepository = (*Adapter)(nil)

// Adapter — адаптер репозитория приглашений в команды (таблица team_invitations).
type Adapter struct {
	invitationReader invitationRepo.InvitationReaderRepository
	invitationWriter invitationRepo.InvitationWriterRepository
}

// NewAdapter создаёт адаптер приглашений.
func NewAdapter(
	invitationReader invitationRepo.InvitationReaderRepository,
	invitationWriter invitationRepo.InvitationWriterRepository,
) *Adapter {
	return &Adapter{
		invitationReader: invitationReader,
		invitationWriter: invitationWriter,
	}
}
