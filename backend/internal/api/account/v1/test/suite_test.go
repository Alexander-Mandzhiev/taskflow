package account_v1_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	account_v1 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/account/v1"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/service/mocks"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

type APISuite struct {
	suite.Suite
	ctx            context.Context // nolint:containedctx
	accountService *mocks.AccountService
	api            *account_v1.API
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.accountService = mocks.NewAccountService(s.T())
	s.api = account_v1.NewAPI(
		s.accountService,
		"access_token",
		15*time.Minute,
		"refresh_token",
		7*24*time.Hour,
		false,
		"",
	)
}

func (s *APISuite) TearDownTest() {}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
