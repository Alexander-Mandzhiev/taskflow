package report

import (
	repo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository"
	reportRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/repository/report"
)

var _ repo.ReportRepository = (*Adapter)(nil)

// Adapter — адаптер репозитория отчётов по задачам.
type Adapter struct {
	reader reportRepo.ReportReaderRepository
}

// NewAdapter создаёт адаптер отчётов.
func NewAdapter(reader reportRepo.ReportReaderRepository) *Adapter {
	return &Adapter{reader: reader}
}
