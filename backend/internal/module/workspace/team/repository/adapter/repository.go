package adapter

import (
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository"
	invitationRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/invitation"
	memberRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/member"
	teamRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/team"
)

var _ repository.TeamRepository = (*Repository)(nil)

// Repository — адаптер поверх team/member/invitation reader/writer.
type Repository struct {
	teamReader       teamRepo.TeamReaderRepository
	teamWriter       teamRepo.TeamWriterRepository
	memberReader     memberRepo.MemberReaderRepository
	memberWriter     memberRepo.MemberWriterRepository
	invitationReader invitationRepo.InvitationReaderRepository
	invitationWriter invitationRepo.InvitationWriterRepository
}

// NewRepository создаёт адаптер.
func NewRepository(
	teamReader teamRepo.TeamReaderRepository,
	teamWriter teamRepo.TeamWriterRepository,
	memberReader memberRepo.MemberReaderRepository,
	memberWriter memberRepo.MemberWriterRepository,
	invitationReader invitationRepo.InvitationReaderRepository,
	invitationWriter invitationRepo.InvitationWriterRepository,
) *Repository {
	return &Repository{
		teamReader:       teamReader,
		teamWriter:       teamWriter,
		memberReader:     memberReader,
		memberWriter:     memberWriter,
		invitationReader: invitationReader,
		invitationWriter: invitationWriter,
	}
}
