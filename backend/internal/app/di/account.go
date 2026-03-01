package di

import (
	"context"
	"fmt"

	account_v1 "mkk/internal/api/account/v1"
	accountRepoDef "mkk/internal/module/identity/account/repository"
	accountRepoCache "mkk/internal/module/identity/account/repository/cache"
	accountServiceDef "mkk/internal/module/identity/account/service"
	accountService "mkk/internal/module/identity/account/service/service"
	"mkk/pkg/password"
)

// AccountV1API возвращает HTTP API account v1 (register, login, logout, whoami).
// Ленивая загрузка сверху вниз: API → Service → Repo → DB. Сначала запрашиваем сервис через геттер, затем собираем API.
func (d *Container) AccountV1API(ctx context.Context) (*account_v1.API, error) {
	if d.accountAPI != nil {
		return d.accountAPI, nil
	}
	svc, err := d.AccountService(ctx)
	if err != nil {
		return nil, fmt.Errorf("account service: %w", err)
	}
	sessionCfg := d.cfg.Session()
	d.accountAPI = account_v1.NewAPI(
		svc,
		sessionCfg.TTL(),
		sessionCfg.IsSecure(),
		sessionCfg.CookieDomain(),
	)
	return d.accountAPI, nil
}

// AccountService возвращает сервис учётных записей (login, register, logout, whoami).
// Ленивая загрузка: запрашиваем репозитории и tx manager через геттеры, затем собираем сервис.
func (d *Container) AccountService(ctx context.Context) (accountServiceDef.AccountService, error) {
	if d.accountService != nil {
		return d.accountService, nil
	}

	sessionRepo, err := d.SessionRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("session repository: %w", err)
	}

	userRepo, err := d.UserRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("user repository: %w", err)
	}

	txMgr, err := d.UserTxManager(ctx)
	if err != nil {
		return nil, fmt.Errorf("user tx manager: %w", err)
	}

	d.accountService = accountService.NewAccountService(
		sessionRepo,
		d.cfg.Session().TTL(),
		userRepo,
		txMgr,
		password.NewBcryptHasher(0),
	)
	return d.accountService, nil
}

// SessionRepository возвращает репозиторий сессий (Redis).
// Ленивая загрузка: запрашиваем Redis через геттер, затем создаём репозиторий.
func (d *Container) SessionRepository(ctx context.Context) (accountRepoDef.SessionCacheRepository, error) {
	if d.sessionRepo != nil {
		return d.sessionRepo, nil
	}
	redisClient, err := d.RedisClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("redis client: %w", err)
	}
	d.sessionRepo = accountRepoCache.NewRepository(redisClient)
	return d.sessionRepo, nil
}
