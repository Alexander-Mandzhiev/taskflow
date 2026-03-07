package task_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	taskRepoMocks "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/mocks"
	svc "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/service"
	taskimpl "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/service/task"
	teamRepoMocks "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/mocks"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/txmanager"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

// ServiceSuite — общий suite для тестов слоя сервиса задач.
// Моки и сервис создаются один раз в SetupSuite; в SetupTest сбрасываются ожидания моков.
type ServiceSuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	taskRepo    *taskRepoMocks.TaskRepository
	historyRepo *taskRepoMocks.TaskHistoryRepository
	teamRepo    *teamRepoMocks.TeamAdapter
	txManager   txmanager.TxManager
	svc         svc.TaskService
}

func (s *ServiceSuite) SetupSuite() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.taskRepo = taskRepoMocks.NewTaskRepository(s.T())
	s.historyRepo = taskRepoMocks.NewTaskHistoryRepository(s.T())
	s.teamRepo = teamRepoMocks.NewTeamAdapter(s.T())
	s.txManager = &txmanager.Noop{}
	s.svc = taskimpl.NewTaskService(s.taskRepo, s.historyRepo, s.teamRepo, s.txManager)
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.taskRepo.ExpectedCalls = nil
	s.historyRepo.ExpectedCalls = nil
	s.teamRepo.ExpectedCalls = nil
}

func TestTaskService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
