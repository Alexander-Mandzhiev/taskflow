package di

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	account_v1 "mkk/internal/api/account/v1"
	accountRepoDef "mkk/internal/module/identity/account/repository"
	accountServiceDef "mkk/internal/module/identity/account/service"
	userRepoDef "mkk/internal/module/identity/user/repository"
	userServiceDef "mkk/internal/module/identity/user/service"
	"mkk/pkg/cache"
	"mkk/pkg/config/contracts"
	"mkk/pkg/database/connectingpool"
	"mkk/pkg/database/txmanager"
)

// Container — DI-контейнер с ленивой инициализацией зависимостей.
// Конфигурация берётся из pkg/config (contracts.Provider). Все методы принимают context.Context и возвращают (dependency, error).
// При первом вызове зависимость создаётся и кешируется; закрытие регистрируется в closer.
type Container struct {
	cfg contracts.Provider

	// Пул БД (MySQL)
	dbPool *connectingpool.Pool
	sqlxDB *sqlx.DB

	// Redis
	redisClient cache.RedisClient
	redisCmdable redis.Cmdable

	// User module (сервисный слой поверх репозиториев)
	userService    userServiceDef.UserService
	userReaderRepo userRepoDef.UserReaderRepository
	userWriterRepo userRepoDef.UserWriterRepository
	userCacheRepo  userRepoDef.UserCacheRepository
	userRepo       userRepoDef.UserRepository
	userTxManager  *txmanager.Manager

	// Account module (session + account service)
	sessionRepo    accountRepoDef.SessionCacheRepository
	accountService accountServiceDef.AccountService
	accountAPI     *account_v1.API
}

// NewContainer создаёт контейнер с конфигурацией из pkg/config (результат config.Load).
func NewContainer(cfg contracts.Provider) *Container {
	return &Container{cfg: cfg}
}
