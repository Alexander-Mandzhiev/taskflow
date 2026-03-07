package task

import (
	taskRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository"
	svc "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/service"
	teamServiceDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/service"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/txmanager"
)

var _ svc.TaskService = (*taskService)(nil)

// taskService — сервис задач и истории изменений.
type taskService struct {
	taskRepo    taskRepo.TaskRepository
	historyRepo taskRepo.TaskHistoryRepository
	teamSvc     teamServiceDef.TeamService
	txManager   txmanager.TxManager
}

// NewTaskService создаёт сервис задач и истории. teamSvc — проверка членства (GetMember) и список команд (ListByUserID).
func NewTaskService(
	taskRepo taskRepo.TaskRepository,
	historyRepo taskRepo.TaskHistoryRepository,
	teamSvc teamServiceDef.TeamService,
	txManager txmanager.TxManager,
) svc.TaskService {
	return &taskService{
		taskRepo:    taskRepo,
		historyRepo: historyRepo,
		teamSvc:     teamSvc,
		txManager:   txManager,
	}
}
