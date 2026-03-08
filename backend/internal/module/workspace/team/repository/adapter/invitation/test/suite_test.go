package invitation_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository"
	invitationAdapter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/adapter/invitation"
	invitationMocks "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/repository/invitation/mocks"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

type AdapterSuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	invitationReader *invitationMocks.InvitationReaderRepository
	invitationWriter *invitationMocks.InvitationWriterRepository
	repo             repository.InvitationRepository
}

func (s *AdapterSuite) SetupTest() {
	s.ctx = context.Background()
	if err := logger.InitDefault(); err != nil {
		panic(err)
	}
	s.invitationReader = invitationMocks.NewInvitationReaderRepository(s.T())
	s.invitationWriter = invitationMocks.NewInvitationWriterRepository(s.T())
	s.repo = invitationAdapter.NewAdapter(s.invitationReader, s.invitationWriter)
}

func TestAdapter(t *testing.T) {
	suite.Run(t, new(AdapterSuite))
}
