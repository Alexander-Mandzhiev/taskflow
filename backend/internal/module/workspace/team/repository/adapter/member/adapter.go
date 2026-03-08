package member

import (
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository"
	memberRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/repository/member"
)

var _ repository.MemberRepository = (*Adapter)(nil)

// Adapter — адаптер репозитория участников команд (таблица team_members).
type Adapter struct {
	memberReader memberRepo.MemberReaderRepository
	memberWriter memberRepo.MemberWriterRepository
}

// NewAdapter создаёт адаптер участников.
func NewAdapter(memberReader memberRepo.MemberReaderRepository, memberWriter memberRepo.MemberWriterRepository) *Adapter {
	return &Adapter{
		memberReader: memberReader,
		memberWriter: memberWriter,
	}
}
