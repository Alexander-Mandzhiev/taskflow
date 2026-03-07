package history_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository"
	historyAdapter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/adapter/history"
	historyMocks "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/repository/history/mocks"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

type AdapterSuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	historyReader *historyMocks.TaskHistoryReaderRepository
	historyWriter *historyMocks.TaskHistoryWriterRepository
	repo          repository.TaskHistoryRepository
}

func (s *AdapterSuite) SetupTest() {
	s.ctx = context.Background()
	if err := logger.InitDefault(); err != nil {
		panic(err)
	}
	s.historyReader = historyMocks.NewTaskHistoryReaderRepository(s.T())
	s.historyWriter = historyMocks.NewTaskHistoryWriterRepository(s.T())
	s.repo = historyAdapter.NewAdapter(s.historyReader, s.historyWriter)
}

func TestAdapter(t *testing.T) {
	suite.Run(t, new(AdapterSuite))
}
