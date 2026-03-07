package team_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	usermocks "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository/mocks"
	teamgrpc "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/client/grpc"
	notificationv1 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/client/grpc/notification/v1"
	teamrepos "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/mocks"
	teamsvc "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/service"
	teamimpl "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/service/service"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/txmanager"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

// ServiceSuite — общий suite для тестов слоя сервиса команд.
// Моки и сервис создаются один раз в SetupSuite; в SetupTest сбрасываются ожидания моков.
type ServiceSuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	teamRepo  *teamrepos.TeamRepository
	userRepo  *usermocks.UserRepository
	txManager txmanager.TxManager
	notifier  teamgrpc.Notification
	svc       teamsvc.TeamService
}

func (s *ServiceSuite) SetupSuite() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.teamRepo = teamrepos.NewTeamRepository(s.T())
	s.userRepo = usermocks.NewUserRepository(s.T())
	s.txManager = &txmanager.Noop{}
	s.notifier = notificationv1.NewClient()
	s.svc = teamimpl.NewTeamService(s.teamRepo, s.txManager, s.userRepo, s.notifier)
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.teamRepo.ExpectedCalls = nil
	s.userRepo.ExpectedCalls = nil
}

func TestTeamService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
