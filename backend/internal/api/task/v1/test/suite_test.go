package task_v1_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	task_v1 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/task/v1"
	taskServiceMocks "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/service/mocks"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

const (
	testTaskID = "660e8400-e29b-41d4-a716-446655440002"
	testUserID = "550e8400-e29b-41d4-a716-446655440000"
)

type APISuite struct {
	suite.Suite
	ctx             context.Context // nolint:containedctx
	taskService     *taskServiceMocks.TaskService
	reportService   *taskServiceMocks.TaskReportService
	commentService  *taskServiceMocks.TaskCommentService
	api             *task_v1.API
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.taskService = taskServiceMocks.NewTaskService(s.T())
	s.reportService = taskServiceMocks.NewTaskReportService(s.T())
	s.commentService = taskServiceMocks.NewTaskCommentService(s.T())
	s.api = task_v1.NewAPI(s.taskService, s.reportService, s.commentService)
}

func (s *APISuite) TearDownTest() {}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
