package team_v1_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	team_v1 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/team/v1"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/service/mocks"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

type APISuite struct {
	suite.Suite
	ctx         context.Context // nolint:containedctx
	teamService *mocks.TeamService
	api         *team_v1.API
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.teamService = mocks.NewTeamService(s.T())
	s.api = team_v1.NewAPI(s.teamService)
}

func (s *APISuite) TearDownTest() {}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
