package adapter_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/repository"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/repository/adapter"
	invitationmocks "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/repository/invitation/mocks"
	membermocks "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/repository/member/mocks"
	teammocks "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/repository/team/mocks"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

type AdapterSuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	teamReader       *teammocks.TeamReaderRepository
	teamWriter       *teammocks.TeamWriterRepository
	memberReader     *membermocks.MemberReaderRepository
	memberWriter     *membermocks.MemberWriterRepository
	invitationReader *invitationmocks.InvitationReaderRepository
	invitationWriter *invitationmocks.InvitationWriterRepository
	repo             repository.TeamRepository
}

func (s *AdapterSuite) SetupTest() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.teamReader = teammocks.NewTeamReaderRepository(s.T())
	s.teamWriter = teammocks.NewTeamWriterRepository(s.T())
	s.memberReader = membermocks.NewMemberReaderRepository(s.T())
	s.memberWriter = membermocks.NewMemberWriterRepository(s.T())
	s.invitationReader = invitationmocks.NewInvitationReaderRepository(s.T())
	s.invitationWriter = invitationmocks.NewInvitationWriterRepository(s.T())
	s.repo = adapter.NewRepository(s.teamReader, s.teamWriter, s.memberReader, s.memberWriter, s.invitationReader, s.invitationWriter)
}

func TestAdapter(t *testing.T) {
	suite.Run(t, new(AdapterSuite))
}
