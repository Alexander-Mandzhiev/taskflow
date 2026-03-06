package service_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"

	accountmocks "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/repository/mocks"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/service"
	svc "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/service/service"
	usermocks "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository/mocks"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/txmanager"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/password"
)

type txManagerStub struct{}

func (t *txManagerStub) WithTx(ctx context.Context, fn func(context.Context, *sqlx.Tx) error, _ ...*sql.TxOptions) error {
	return fn(ctx, nil)
}

func (t *txManagerStub) WithSerializableTx(ctx context.Context, fn func(context.Context, *sqlx.Tx) error) error {
	return fn(ctx, nil)
}

var _ txmanager.TxManager = (*txManagerStub)(nil)

type ServiceSuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	sessionRepo *accountmocks.SessionCacheRepository
	userRepo    *usermocks.UserRepository
	txManager   *txManagerStub
	hasher      password.Hasher
	svc         service.AccountService
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.sessionRepo = accountmocks.NewSessionCacheRepository(s.T())
	s.userRepo = usermocks.NewUserRepository(s.T())
	s.txManager = &txManagerStub{}
	s.hasher = password.NewBcryptHasher(4)
	s.svc = svc.NewAccountService(
		s.sessionRepo,
		24*time.Hour,
		s.userRepo,
		s.txManager,
		s.hasher,
	)
}

func (s *ServiceSuite) TearDownTest() {}

func TestAccountService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
