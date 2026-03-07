package service

import (
	userRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository"
	teamRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/repository"
	def "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/service"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/txmanager"
)

var _ def.TeamService = (*teamService)(nil)

type teamService struct {
	repo      teamRepo.TeamRepository
	txManager txmanager.TxManager
	userRepo  userRepo.UserRepository
}

// NewTeamService создаёт сервис команд. userRepo — адаптер пользователей (для invite по email).
func NewTeamService(repo teamRepo.TeamRepository, txManager txmanager.TxManager, userRepo userRepo.UserRepository) def.TeamService {
	return &teamService{
		repo:      repo,
		txManager: txManager,
		userRepo:  userRepo,
	}
}
