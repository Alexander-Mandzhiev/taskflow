package di

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	account_v1 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/account/v1"
	team_v1 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/team/v1"
	accountRepoDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/repository"
	accountServiceDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/service"
	userRepoAdapterDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository"
	userCacheRepoDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository/cache"
	userReaderRepoDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository/user"
	userServiceDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/service"
	teamRepoDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository"
	teamServiceDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/service"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/cache"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/closer"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/contracts"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/connectingpool"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/txmanager"
)

// Container — DI-контейнер с ленивой инициализацией зависимостей.
// Конфигурация берётся из pkg/config (contracts.Provider). Все методы принимают context.Context и возвращают (dependency, error).
// При первом вызове зависимость создаётся и кешируется; закрытие регистрируется в переданном closer (SetCloser до первого использования).
type Container struct {
	cfg contracts.Provider
	cl  *closer.Closer

	// Пул БД (MySQL)
	dbPool *connectingpool.Pool
	sqlxDB *sqlx.DB

	// Redis
	redisClient  cache.RedisClient
	redisCmdable redis.Cmdable

	// User module (сервисный слой поверх репозиториев)
	userService    userServiceDef.UserService
	userReaderRepo userReaderRepoDef.UserReaderRepository
	userWriterRepo userReaderRepoDef.UserWriterRepository
	userCacheRepo  userCacheRepoDef.UserCacheRepository
	userRepo       userRepoAdapterDef.UserRepository
	userTxManager  *txmanager.Manager

	// Account module (session + account service)
	sessionRepo        accountRepoDef.SessionCacheRepository
	accountService     accountServiceDef.AccountService
	accountAPI         *account_v1.API
	accountMiddlewares *accountMiddlewares

	// Team module
	teamRepo    teamRepoDef.TeamRepository
	teamService teamServiceDef.TeamService
	teamAPI     *team_v1.API
}

// NewContainer создаёт контейнер с конфигурацией из pkg/config (результат config.Load).
// Перед вызовом SqlxDB/RedisClient/RegisterAccountRoutes необходимо вызвать SetCloser.
func NewContainer(cfg contracts.Provider) *Container {
	return &Container{cfg: cfg}
}

// SetCloser задаёт менеджер ресурсов для регистрации закрытия (вызывается из App.initCloser).
func (d *Container) SetCloser(c *closer.Closer) {
	d.cl = c
}

func (d *Container) requireCloser() error {
	if d.cl == nil {
		return closer.ErrNotSet
	}
	return nil
}
