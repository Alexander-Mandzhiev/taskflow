package team_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository"
	teamAdapter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/adapter/team"
	teamMocks "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/repository/team/mocks"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

type AdapterSuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	teamReader *teamMocks.TeamReaderRepository
	teamWriter *teamMocks.TeamWriterRepository
	repo       repository.TeamRepository
}

func (s *AdapterSuite) SetupTest() {
	s.ctx = context.Background()
	if err := logger.InitDefault(); err != nil {
		panic(err)
	}
	s.teamReader = teamMocks.NewTeamReaderRepository(s.T())
	s.teamWriter = teamMocks.NewTeamWriterRepository(s.T())
	s.repo = teamAdapter.NewAdapter(s.teamReader, s.teamWriter)
}

func TestAdapter(t *testing.T) {
	suite.Run(t, new(AdapterSuite))
}
