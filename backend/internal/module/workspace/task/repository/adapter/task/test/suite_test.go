package task_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository"
	taskAdapter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/adapter/task"
	taskMocks "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/repository/task/mocks"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

type AdapterSuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	taskReader *taskMocks.TaskReaderRepository
	taskWriter *taskMocks.TaskWriterRepository
	repo       repository.TaskRepository
}

func (s *AdapterSuite) SetupTest() {
	s.ctx = context.Background()
	if err := logger.InitDefault(); err != nil {
		panic(err)
	}
	s.taskReader = taskMocks.NewTaskReaderRepository(s.T())
	s.taskWriter = taskMocks.NewTaskWriterRepository(s.T())
	s.repo = taskAdapter.NewAdapter(s.taskReader, s.taskWriter, nil)
}

func TestAdapter(t *testing.T) {
	suite.Run(t, new(AdapterSuite))
}
