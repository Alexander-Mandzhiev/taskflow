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

	// JWT: для генерации access и refresh токенов в Login
	accessSecret  string
	refreshSecret string
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

// NewAccountService создаёт сервис account.
// accessSecret, refreshSecret, accessTTL, refreshTTL — параметры для генерации JWT в Login.
func NewAccountService(
	sessionRepo repository.SessionCacheRepository,
	sessionTTL time.Duration,
	userRepo userrepo.UserRepository,
	txManager txmanager.TxManager,
	hasher password.Hasher,
	accessSecret string,
	refreshSecret string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) def.AccountService {
	return &accountService{
		sessionRepo:   sessionRepo,
		sessionTTL:    sessionTTL,
		userRepo:      userRepo,
		txManager:     txManager,
		hasher:        hasher,
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}
