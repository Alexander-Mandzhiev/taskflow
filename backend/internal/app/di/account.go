package di

import (
	"context"
	"fmt"

	account_v1 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/account/v1"
	accountRepoDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/repository"
	accountRepoCache "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/repository/cache"
	accountServiceDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/service"
	accountService "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/service/service"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/password"
)

// AccountV1API возвращает HTTP API account v1 (register, login, logout).
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
	jwtCfg := d.cfg.JWT()
	d.accountAPI = account_v1.NewAPI(
		svc,
		jwtCfg.AccessTokenCookieName(),
		jwtCfg.AccessTTL(),
		jwtCfg.RefreshTokenCookieName(),
		jwtCfg.RefreshTTL(),
		sessionCfg.IsSecure(),
		sessionCfg.CookieDomain(),
	)
	return d.accountAPI, nil
}

// AccountService возвращает сервис учётных записей (login, register, logout).
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

	jwtCfg := d.cfg.JWT()
	d.accountService = accountService.NewAccountService(
		sessionRepo,
		jwtCfg.RefreshTTL(), // TTL сессии в Redis = время жизни refresh-токена
		userRepo,
		txMgr,
		password.NewBcryptHasher(0),
		jwtCfg.AccessSecret(),
		jwtCfg.RefreshSecret(),
		jwtCfg.AccessTTL(),
		jwtCfg.RefreshTTL(),
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
