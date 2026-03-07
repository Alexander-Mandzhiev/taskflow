package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// TaskService — слой сервиса задач и истории.
// Транзакции открываются внутри сервиса (txmanager.WithTx). ID передаются как uuid.UUID.
// Права: создание/чтение/обновление/история — только для членов команды задачи.
// Отчёты (TeamTaskStats, TopCreatorsByTeam, TasksWithInvalidAssignee) реализация сервиса получает через ReportRepository из модуля workspace/report.
//
// Валидация: Create — статус проверяется в репозитории (конвертер); при заданном AssigneeID сервис должен проверить через TeamService.GetMember(teamID, assigneeID), при ErrMemberNotFound возвращать model.ErrAssigneeNotInTeam.
// Update — статус через model.ValidateTaskInput; assignee — так же через GetMember(teamID, assigneeID).
type TaskService interface {
	// Create создаёт задачу. userID — участник команды; teamID передаётся в сигнатуре. При nil input — model.ErrNilInput. При заданном AssigneeID — должен быть в команде, иначе model.ErrAssigneeNotInTeam.
	Create(ctx context.Context, userID, teamID uuid.UUID, input *model.TaskInput) (*model.Task, error)

	// GetByID возвращает задачу только если userID — участник команды задачи; иначе model.ErrForbidden или model.ErrTaskNotFound.
	GetByID(ctx context.Context, taskID, userID uuid.UUID) (*model.Task, error)

	// List возвращает список задач с фильтром и пагинацией. Если задан filter.TeamID — userID должен быть в этой команде.
	List(ctx context.Context, userID uuid.UUID, filter *model.TaskListFilter, pagination *model.TaskPagination) ([]*model.Task, int, error)

	// Update полностью обновляет задачу (title, description, status, assignee_id). taskID в сигнатуре. userID — участник команды; при изменениях пишется task_history. Статус — model.ValidateTaskInput; при заданном AssigneeID — в команде, иначе model.ErrAssigneeNotInTeam.
	Update(ctx context.Context, userID, taskID uuid.UUID, input *model.TaskInput) (*model.Task, error)

	// Delete выполняет мягкое удаление задачи (deleted_at). userID должен быть участником команды задачи; иначе model.ErrForbidden или model.ErrTaskNotFound.
	Delete(ctx context.Context, userID, taskID uuid.UUID) error

	// Restore снимает пометку удаления. userID должен быть участником команды задачи; иначе model.ErrForbidden или model.ErrTaskNotFound.
	Restore(ctx context.Context, userID, taskID uuid.UUID) (*model.Task, error)

	// GetHistory возвращает историю изменений задачи. userID должен быть участником команды задачи.
	GetHistory(ctx context.Context, taskID, userID uuid.UUID) ([]*model.TaskHistory, error)

	// TeamTaskStats — отчёт (а): по каждой команде, где userID участник — название, кол-во участников, кол-во задач done за since.
	TeamTaskStats(ctx context.Context, userID uuid.UUID, since time.Time) ([]*model.TeamTaskStats, error)

	// TopCreatorsByTeam — отчёт (б): топ-N по созданным задачам в каждой команде за период (только команды userID).
	TopCreatorsByTeam(ctx context.Context, userID uuid.UUID, since time.Time, limit int) ([]*model.TeamTopCreator, error)

	// TasksWithInvalidAssignee — отчёт (в): задачи с assignee не из команды (только по командам userID).
	TasksWithInvalidAssignee(ctx context.Context, userID uuid.UUID) ([]*model.Task, error)
}
