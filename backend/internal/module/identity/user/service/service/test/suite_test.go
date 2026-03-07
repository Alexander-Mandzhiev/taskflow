package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository/mocks"
	userservice "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/service"
	svc "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/service/service"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/txmanager"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

// ServiceSuite — общий suite для тестов слоя сервиса пользователей.
// Моки и сервис создаются один раз в SetupSuite; в SetupTest сбрасываются ожидания моков.
type ServiceSuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	repo      *mocks.UserRepository
	txManager txmanager.TxManager
	svc       userservice.UserService
}

func (s *ServiceSuite) SetupSuite() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.repo = mocks.NewUserRepository(s.T())
	s.txManager = &txmanager.Noop{}
	s.svc = svc.NewUserService(s.repo, s.txManager)
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.repo.ExpectedCalls = nil
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
