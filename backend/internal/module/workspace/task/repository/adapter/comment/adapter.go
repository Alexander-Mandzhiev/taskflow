package comment

import (
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository"
	commentRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/repository/comment"
)

var _ repository.TaskCommentRepository = (*Adapter)(nil)

// Adapter — адаптер репозитория комментариев к задачам (таблица task_comments).
type Adapter struct {
	commentReader commentRepo.TaskCommentReaderRepository
	commentWriter commentRepo.TaskCommentWriterRepository
}

// NewAdapter создаёт адаптер комментариев к задачам.
func NewAdapter(
	commentReader commentRepo.TaskCommentReaderRepository,
	commentWriter commentRepo.TaskCommentWriterRepository,
) *Adapter {
	return &Adapter{
		commentReader: commentReader,
		commentWriter: commentWriter,
	}
}
