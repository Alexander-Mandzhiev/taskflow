package team

import (
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository"
	teamRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/repository/team"
)

var _ repository.TeamRepository = (*Adapter)(nil)

// Adapter — адаптер репозитория команд (таблица teams).
type Adapter struct {
	teamReader teamRepo.TeamReaderRepository
	teamWriter teamRepo.TeamWriterRepository
}

// NewAdapter создаёт адаптер команд.
func NewAdapter(teamReader teamRepo.TeamReaderRepository, teamWriter teamRepo.TeamWriterRepository) *Adapter {
	return &Adapter{
		teamReader: teamReader,
		teamWriter: teamWriter,
	}
}
