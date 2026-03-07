package task

import (
	taskRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository"
	svc "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/service"
	teamRepoDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/txmanager"
)

var _ svc.TaskService = (*taskService)(nil)

// taskService — сервис задач и истории изменений.
// teamRepo — проверка членства (GetMember) в той же транзакции, что и запись (Create/Update).
type taskService struct {
	taskRepo    taskRepo.TaskRepository
	historyRepo taskRepo.TaskHistoryRepository
	teamRepo    teamRepoDef.TeamAdapter
	txManager   txmanager.TxManager
}

// NewTaskService создаёт сервис задач. teamRepo — адаптер команд для GetMember(ctx, tx, ...); передаём tx, чтобы проверка и запись были в одной транзакции.
func NewTaskService(
	taskRepo taskRepo.TaskRepository,
	historyRepo taskRepo.TaskHistoryRepository,
	teamRepo teamRepoDef.TeamAdapter,
	txManager txmanager.TxManager,
) svc.TaskService {
	return &taskService{
		taskRepo:    taskRepo,
		historyRepo: historyRepo,
		teamRepo:    teamRepo,
		txManager:   txManager,
	}
}
