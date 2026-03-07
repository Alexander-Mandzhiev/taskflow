package task

import (
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository"
	taskRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/repository/task"
)

var _ repository.TaskRepository = (*Adapter)(nil)

// Adapter — адаптер репозитория задач (таблица tasks).
type Adapter struct {
	taskReader taskRepo.TaskReaderRepository
	taskWriter taskRepo.TaskWriterRepository
}

// NewAdapter создаёт адаптер задач.
func NewAdapter(
	taskReader taskRepo.TaskReaderRepository,
	taskWriter taskRepo.TaskWriterRepository,
) *Adapter {
	return &Adapter{
		taskReader: taskReader,
		taskWriter: taskWriter,
	}
}
