package adapter

import (
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository"
	invitationRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/invitation"
	memberRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/member"
	teamRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/team"
)

var _ repository.TeamAdapter = (*Adapter)(nil)

// Adapter — адаптер поверх team/member/invitation reader/writer.
type Adapter struct {
	teamReader       teamRepo.TeamReaderRepository
	teamWriter       teamRepo.TeamWriterRepository
	memberReader     memberRepo.MemberReaderRepository
	memberWriter     memberRepo.MemberWriterRepository
	invitationReader invitationRepo.InvitationReaderRepository
	invitationWriter invitationRepo.InvitationWriterRepository
}

// NewAdapter создаёт адаптер.
func NewAdapter(
	teamReader teamRepo.TeamReaderRepository,
	teamWriter teamRepo.TeamWriterRepository,
	memberReader memberRepo.MemberReaderRepository,
	memberWriter memberRepo.MemberWriterRepository,
	invitationReader invitationRepo.InvitationReaderRepository,
	invitationWriter invitationRepo.InvitationWriterRepository,
) *Adapter {
	return &Adapter{
		teamReader:       teamReader,
		teamWriter:       teamWriter,
		memberReader:     memberReader,
		memberWriter:     memberWriter,
		invitationReader: invitationReader,
		invitationWriter: invitationWriter,
	}
}
