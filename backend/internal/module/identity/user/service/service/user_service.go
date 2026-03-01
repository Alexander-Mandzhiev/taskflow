package service

import (
	"mkk/internal/module/identity/user/repository"
	api "mkk/internal/module/identity/user/service"
	"mkk/pkg/database/txmanager"
)

var _ api.UserService = (*userService)(nil)

type userService struct {
	repo      repository.UserRepository
	txManager txmanager.TxManager
}

// NewUserService создаёт сервис пользователей с заданным репозиторием и менеджером транзакций.
func NewUserService(repo repository.UserRepository, txManager txmanager.TxManager) api.UserService {
	return &userService{
		repo:      repo,
		txManager: txManager,
	}
}
