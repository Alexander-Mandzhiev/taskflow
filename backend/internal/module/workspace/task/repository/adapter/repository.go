package adapter

import (
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository"
	historyRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/history"
	taskRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/task"
)

var _ repository.TaskRepository = (*Repository)(nil)

// Repository — адаптер поверх task reader/writer и history reader/writer. Отчёты — отдельный модуль workspace/report.
type Repository struct {
	taskReader    taskRepo.TaskReaderRepository
	taskWriter    taskRepo.TaskWriterRepository
	historyReader historyRepo.TaskHistoryReaderRepository
	historyWriter historyRepo.TaskHistoryWriterRepository
}

// NewRepository создаёт адаптер.
func NewRepository(
	taskReader taskRepo.TaskReaderRepository,
	taskWriter taskRepo.TaskWriterRepository,
	historyReader historyRepo.TaskHistoryReaderRepository,
	historyWriter historyRepo.TaskHistoryWriterRepository,
) *Repository {
	return &Repository{
		taskReader:    taskReader,
		taskWriter:    taskWriter,
		historyReader: historyReader,
		historyWriter: historyWriter,
	}
}
