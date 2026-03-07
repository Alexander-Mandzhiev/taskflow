package adapter_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/adapter"
	mocks2 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/invitation/mocks"
	mocks3 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/member/mocks"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/team/mocks"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

type AdapterSuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	teamReader       *mocks.TeamReaderRepository
	teamWriter       *mocks.TeamWriterRepository
	memberReader     *mocks3.MemberReaderRepository
	memberWriter     *mocks3.MemberWriterRepository
	invitationReader *mocks2.InvitationReaderRepository
	invitationWriter *mocks2.InvitationWriterRepository
	repo             repository.TeamRepository
}

func (s *AdapterSuite) SetupTest() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.teamReader = mocks.NewTeamReaderRepository(s.T())
	s.teamWriter = mocks.NewTeamWriterRepository(s.T())
	s.memberReader = mocks3.NewMemberReaderRepository(s.T())
	s.memberWriter = mocks3.NewMemberWriterRepository(s.T())
	s.invitationReader = mocks2.NewInvitationReaderRepository(s.T())
	s.invitationWriter = mocks2.NewInvitationWriterRepository(s.T())
	s.repo = adapter.NewRepository(s.teamReader, s.teamWriter, s.memberReader, s.memberWriter, s.invitationReader, s.invitationWriter)
}

func TestAdapter(t *testing.T) {
	suite.Run(t, new(AdapterSuite))
}
