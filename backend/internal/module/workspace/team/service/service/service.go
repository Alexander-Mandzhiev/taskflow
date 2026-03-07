package team

import (
	userRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository"
	teamClient "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/client/grpc"
	teamRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository"
	def "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/service"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/txmanager"
)

var _ def.TeamService = (*teamService)(nil)

type teamService struct {
	repo      teamRepo.TeamAdapter
	txManager txmanager.TxManager
	userRepo  userRepo.UserRepository
	notifier  teamClient.Notification
}

// NewTeamService создаёт сервис команд. userRepo — для invite по email; notifier — отправка уведомления (мок или gRPC). Ссылку «принять приглашение» собирает сервис уведомлений из inv.Token.
func NewTeamService(repo teamRepo.TeamAdapter, txManager txmanager.TxManager, userRepo userRepo.UserRepository, notifier teamClient.Notification) def.TeamService {
	return &teamService{
		repo:      repo,
		txManager: txManager,
		userRepo:  userRepo,
		notifier:  notifier,
	}
}
