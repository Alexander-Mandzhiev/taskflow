package di

import (
	"context"
	"fmt"

	userRepoDef "mkk/internal/module/identity/user/repository"
	userRepoAdapter "mkk/internal/module/identity/user/repository/adapter"
	userRepoCache "mkk/internal/module/identity/user/repository/cache"
	userRepoReader "mkk/internal/module/identity/user/repository/reader"
	userRepoWriter "mkk/internal/module/identity/user/repository/writer"
	userServiceDef "mkk/internal/module/identity/user/service"
	userServiceImpl "mkk/internal/module/identity/user/service/service"
	"mkk/pkg/database/txmanager"
)

// UserService возвращает сервисный слой пользователей (CRUD, профиль, смена пароля).
// Ленивая загрузка: запрашиваем UserRepository и UserTxManager через геттеры, затем собираем сервис.
func (d *Container) UserService(ctx context.Context) (userServiceDef.UserService, error) {
	if d.userService != nil {
		return d.userService, nil
	}
	userRepo, err := d.UserRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("user repository: %w", err)
	}
	txMgr, err := d.UserTxManager(ctx)
	if err != nil {
		return nil, fmt.Errorf("user tx manager: %w", err)
	}
	d.userService = userServiceImpl.NewUserService(userRepo, txMgr)
	return d.userService, nil
}

// UserRepository возвращает адаптер репозитория (reader + writer + cache).
// Ленивая загрузка: запрашиваем reader, writer, cache через геттеры, затем собираем адаптер.
func (d *Container) UserRepository(ctx context.Context) (userRepoDef.UserRepository, error) {
	if d.userRepo != nil {
		return d.userRepo, nil
	}

	reader, err := d.UserReaderRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("user reader: %w", err)
	}

	writer, err := d.UserWriterRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("user writer: %w", err)
	}

	cacheRepo, err := d.UserCacheRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("user cache: %w", err)
	}

	d.userRepo = userRepoAdapter.NewRepository(reader, writer, cacheRepo)

	return d.userRepo, nil
}

// UserTxManager возвращает менеджер транзакций для модуля user.
// Ленивая загрузка: запрашиваем SqlxDB через геттер, затем создаём менеджер.
func (d *Container) UserTxManager(ctx context.Context) (*txmanager.Manager, error) {
	if d.userTxManager != nil {
		return d.userTxManager, nil
	}

	db, err := d.SqlxDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("sqlx db: %w", err)
	}

	d.userTxManager = txmanager.New(db, "user.txmanager")

	return d.userTxManager, nil
}

// UserReaderRepository возвращает репозиторий чтения пользователей.
// Ленивая загрузка: запрашиваем SqlxDB через геттер.
func (d *Container) UserReaderRepository(ctx context.Context) (userRepoDef.UserReaderRepository, error) {
	if d.userReaderRepo != nil {
		return d.userReaderRepo, nil
	}

	db, err := d.SqlxDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("sqlx db: %w", err)
	}

	d.userReaderRepo = userRepoReader.NewRepository(db)

	return d.userReaderRepo, nil
}

// UserWriterRepository возвращает репозиторий записи пользователей.
// Ленивая загрузка: запрашиваем SqlxDB через геттер.
func (d *Container) UserWriterRepository(ctx context.Context) (userRepoDef.UserWriterRepository, error) {
	if d.userWriterRepo != nil {
		return d.userWriterRepo, nil
	}

	db, err := d.SqlxDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("sqlx db: %w", err)
	}

	d.userWriterRepo = userRepoWriter.NewRepository(db)

	return d.userWriterRepo, nil
}

// UserCacheRepository возвращает кеш-репозиторий пользователей.
// Ленивая загрузка: запрашиваем RedisClient через геттер.
func (d *Container) UserCacheRepository(ctx context.Context) (userRepoDef.UserCacheRepository, error) {
	if d.userCacheRepo != nil {
		return d.userCacheRepo, nil
	}

	redisClient, err := d.RedisClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("redis client: %w", err)
	}

	d.userCacheRepo = userRepoCache.NewRepository(redisClient)

	return d.userCacheRepo, nil
}
