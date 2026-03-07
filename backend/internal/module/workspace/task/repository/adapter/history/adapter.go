package history

import (
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository"
	historyRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/repository/history"
)

var _ repository.TaskHistoryRepository = (*Adapter)(nil)

// Adapter — адаптер репозитория истории задач (таблица task_history).
type Adapter struct {
	historyReader historyRepo.TaskHistoryReaderRepository
	historyWriter historyRepo.TaskHistoryWriterRepository
}

// NewAdapter создаёт адаптер истории задач.
func NewAdapter(
	historyReader historyRepo.TaskHistoryReaderRepository,
	historyWriter historyRepo.TaskHistoryWriterRepository,
) *Adapter {
	return &Adapter{
		historyReader: historyReader,
		historyWriter: historyWriter,
	}
}
