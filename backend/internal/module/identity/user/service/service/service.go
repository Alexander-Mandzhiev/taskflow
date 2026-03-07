package user

import (
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository"
	def "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/service"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/txmanager"
)

var _ def.UserService = (*userService)(nil)

type userService struct {
	repo      repository.UserRepository
	txManager txmanager.TxManager
}

// NewUserService создаёт сервис пользователей с заданным репозиторием и менеджером транзакций.
func NewUserService(repo repository.UserRepository, txManager txmanager.TxManager) def.UserService {
	return &userService{
		repo:      repo,
		txManager: txManager,
	}
}
