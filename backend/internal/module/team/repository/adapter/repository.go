package adapter

import (
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/repository"
	memberRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/repository/member"
	teamRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/repository/team"
)

var _ repository.TeamRepository = (*Repository)(nil)

// Repository — адаптер поверх TeamReaderRepository, TeamWriterRepository, MemberReaderRepository, MemberWriterRepository.
type Repository struct {
	teamReader   teamRepo.TeamReaderRepository
	teamWriter   teamRepo.TeamWriterRepository
	memberReader memberRepo.MemberReaderRepository
	memberWriter memberRepo.MemberWriterRepository
}

// NewRepository создаёт адаптер.
func NewRepository(
	teamReader teamRepo.TeamReaderRepository,
	teamWriter teamRepo.TeamWriterRepository,
	memberReader memberRepo.MemberReaderRepository,
	memberWriter memberRepo.MemberWriterRepository,
) *Repository {
	return &Repository{
		teamReader:   teamReader,
		teamWriter:   teamWriter,
		memberReader: memberReader,
		memberWriter: memberWriter,
	}
}
