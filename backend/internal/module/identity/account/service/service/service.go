package service

import (
	"time"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/repository"
	def "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/service"
	userrepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/txmanager"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/password"
)

var _ def.AccountService = (*accountService)(nil)

type accountService struct {
	sessionRepo repository.SessionCacheRepository

	sessionTTL time.Duration

	userRepo userrepo.UserRepository

	txManager txmanager.TxManager

	hasher password.Hasher
}

// NewAccountService создаёт сервис account.

func NewAccountService(
	sessionRepo repository.SessionCacheRepository,

	sessionTTL time.Duration,

	userRepo userrepo.UserRepository,

	txManager txmanager.TxManager,

	hasher password.Hasher,
) def.AccountService {
	return &accountService{
		sessionRepo: sessionRepo,

		sessionTTL: sessionTTL,

		userRepo: userRepo,

		txManager: txManager,

		hasher: hasher,
	}
}
