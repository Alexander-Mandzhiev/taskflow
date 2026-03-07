package comment_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	taskRepoMocks "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/mocks"
	svc "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/service"
	commentImpl "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/service/comment"
	teamRepoMocks "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/mocks"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/txmanager"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

// ServiceSuite — общий suite для тестов сервиса комментариев.
// Моки и сервис создаются в SetupSuite; в SetupTest сбрасываются ожидания моков.
type ServiceSuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	taskRepo    *taskRepoMocks.TaskRepository
	commentRepo *taskRepoMocks.TaskCommentRepository
	teamRepo    *teamRepoMocks.TeamAdapter
	txManager   txmanager.TxManager
	svc         svc.TaskCommentService
}

func (s *ServiceSuite) SetupSuite() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.taskRepo = taskRepoMocks.NewTaskRepository(s.T())
	s.commentRepo = taskRepoMocks.NewTaskCommentRepository(s.T())
	s.teamRepo = teamRepoMocks.NewTeamAdapter(s.T())
	s.txManager = &txmanager.Noop{}
	s.svc = commentImpl.NewService(s.taskRepo, s.commentRepo, s.teamRepo, s.txManager)
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.taskRepo.ExpectedCalls = nil
	s.commentRepo.ExpectedCalls = nil
	s.teamRepo.ExpectedCalls = nil
}

func TestCommentService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
