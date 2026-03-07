package report_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	taskRepoMocks "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/mocks"
	svc "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/service"
	reportimpl "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/service/report"
	teamSvcMocks "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/service/mocks"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

// ServiceSuite — общий suite для тестов слоя сервиса отчётов по задачам.
// Моки и сервис создаются один раз в SetupSuite; в SetupTest сбрасываются ожидания моков.
type ServiceSuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	reportRepo *taskRepoMocks.ReportRepository
	teamSvc    *teamSvcMocks.TeamService
	svc        svc.TaskReportService
}

func (s *ServiceSuite) SetupSuite() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.reportRepo = taskRepoMocks.NewReportRepository(s.T())
	s.teamSvc = teamSvcMocks.NewTeamService(s.T())
	s.svc = reportimpl.NewTaskReportService(s.reportRepo, s.teamSvc)
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.reportRepo.ExpectedCalls = nil
	s.teamSvc.ExpectedCalls = nil
}

func TestReportService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
