package service

import (
	"time"

	def "mkk/internal/module/identity/account/service"
	"mkk/internal/module/identity/account/repository"
	userrepo "mkk/internal/module/identity/user/repository"
	"mkk/pkg/database/txmanager"
	"mkk/pkg/password"
)

var _ def.AccountService = (*accountService)(nil)

type accountService struct {
	sessionRepo repository.SessionCacheRepository
	sessionTTL  time.Duration
	userRepo    userrepo.UserRepository
	txManager   txmanager.TxManager
	hasher      password.Hasher
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
		sessionTTL:  sessionTTL,
		userRepo:    userRepo,
		txManager:   txManager,
		hasher:      hasher,
	}
}
