package member_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository"
	memberAdapter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/adapter/member"
	memberMocks "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/repository/member/mocks"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

type AdapterSuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	memberReader *memberMocks.MemberReaderRepository
	memberWriter *memberMocks.MemberWriterRepository
	repo         repository.MemberRepository
}

func (s *AdapterSuite) SetupTest() {
	s.ctx = context.Background()
	if err := logger.InitDefault(); err != nil {
		panic(err)
	}
	s.memberReader = memberMocks.NewMemberReaderRepository(s.T())
	s.memberWriter = memberMocks.NewMemberWriterRepository(s.T())
	s.repo = memberAdapter.NewAdapter(s.memberReader, s.memberWriter)
}

func TestAdapter(t *testing.T) {
	suite.Run(t, new(AdapterSuite))
}
