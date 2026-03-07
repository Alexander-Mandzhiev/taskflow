package adapter_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository/adapter"
	cachemocks "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository/cache/mocks"
	usermocks "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository/user/mocks"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

type AdapterSuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	reader *usermocks.UserReaderRepository
	writer *usermocks.UserWriterRepository
	cache  *cachemocks.UserCacheRepository
	repo   repository.UserRepository
}

func (s *AdapterSuite) SetupTest() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.reader = usermocks.NewUserReaderRepository(s.T())
	s.writer = usermocks.NewUserWriterRepository(s.T())
	s.cache = cachemocks.NewUserCacheRepository(s.T())
	s.repo = adapter.NewAdapter(s.reader, s.writer, s.cache)
}

func TestAdapter(t *testing.T) {
	suite.Run(t, new(AdapterSuite))
}
