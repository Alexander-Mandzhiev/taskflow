package comment

import (
	taskRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/service"
	teamRepoDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/txmanager"
)

var _ service.TaskCommentService = (*commentService)(nil)

// commentService — сервис комментариев к задачам. Доступ только для участников команды задачи.
type commentService struct {
	taskRepo    taskRepo.TaskRepository
	commentRepo taskRepo.TaskCommentRepository
	memberRepo  teamRepoDef.MemberRepository
	txManager   txmanager.TxManager
}

// NewCommentService создаёт сервис комментариев. taskRepo и memberRepo — для проверки доступа (user в команде задачи).
func NewCommentService(
	taskRepo taskRepo.TaskRepository,
	commentRepo taskRepo.TaskCommentRepository,
	memberRepo teamRepoDef.MemberRepository,
	txManager txmanager.TxManager,
) service.TaskCommentService {
	return &commentService{
		taskRepo:    taskRepo,
		commentRepo: commentRepo,
		memberRepo:  memberRepo,
		txManager:   txManager,
	}
}
