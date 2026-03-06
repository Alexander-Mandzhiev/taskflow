package service_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository/mocks"
	userservice "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/service"
	svc "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/service/service"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/txmanager"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

type txManagerStub struct {
	tx          *sqlx.Tx
	withTxCalls int
	forceErr    error
}

func (t *txManagerStub) WithTx(ctx context.Context, fn func(context.Context, *sqlx.Tx) error, _ ...*sql.TxOptions) error {
	t.withTxCalls++
	if t.forceErr != nil {
		return t.forceErr
	}
	if t.tx == nil {
		t.tx = &sqlx.Tx{}
	}
	return fn(ctx, t.tx)
}

func (t *txManagerStub) WithSerializableTx(ctx context.Context, fn func(context.Context, *sqlx.Tx) error) error {
	if t.forceErr != nil {
		return t.forceErr
	}
	if t.tx == nil {
		t.tx = &sqlx.Tx{}
	}
	return fn(ctx, t.tx)
}

var _ txmanager.TxManager = (*txManagerStub)(nil)

type ServiceSuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	repo      *mocks.UserRepository
	txManager *txManagerStub
	svc       userservice.UserService
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.repo = mocks.NewUserRepository(s.T())
	s.txManager = &txManagerStub{}
	s.svc = svc.NewUserService(s.repo, s.txManager)
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
