package comment

import (
	taskRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/service"
	teamRepoDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/txmanager"
)

var _ service.TaskCommentService = (*Service)(nil)

// Service — сервис комментариев к задачам. Доступ только для участников команды задачи.
type Service struct {
	taskRepo    taskRepo.TaskRepository
	commentRepo taskRepo.TaskCommentRepository
	teamRepo    teamRepoDef.TeamAdapter
	txManager   txmanager.TxManager
}

// NewService создаёт сервис комментариев. taskRepo и teamRepo — для проверки доступа (user в команде задачи).
func NewService(
	taskRepo taskRepo.TaskRepository,
	commentRepo taskRepo.TaskCommentRepository,
	teamRepo teamRepoDef.TeamAdapter,
	txManager txmanager.TxManager,
) *Service {
	return &Service{
		taskRepo:    taskRepo,
		commentRepo: commentRepo,
		teamRepo:    teamRepo,
		txManager:   txManager,
	}
}
