package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

// TaskService — CRUD задач, мягкое удаление/восстановление, история изменений.
// Транзакции открываются внутри сервиса (txmanager.WithTx). ID передаются как uuid.UUID.
// Права: создание/чтение/обновление/история — только для членов команды задачи.
//
// Валидация: Create — статус в репозитории (конвертер); при заданном AssigneeID — TeamService.GetMember(teamID, assigneeID), при ErrMemberNotFound → model.ErrAssigneeNotInTeam.
// Update — статус через model.ValidateTaskInput; assignee — так же через GetMember.
type TaskService interface {
	// Create создаёт задачу. userID — участник команды; teamID в сигнатуре. При nil input — model.ErrNilInput. AssigneeID — должен быть в команде.
	Create(ctx context.Context, userID, teamID uuid.UUID, input *model.TaskInput) (*model.Task, error)

	// GetByID возвращает задачу только если userID — участник команды; иначе model.ErrForbidden или model.ErrTaskNotFound.
	GetByID(ctx context.Context, taskID, userID uuid.UUID) (*model.Task, error)

	// List — список задач по фильтру (критерии + limit/offset в filter). Валидация filter (limit > 0) в сервисе/API. filter.TeamID обязателен, userID — в команде.
	List(ctx context.Context, userID uuid.UUID, filter *model.TaskListFilter) ([]*model.Task, int, error)

	// Update полностью обновляет задачу (title, description, status, assignee_id). При изменениях пишется task_history.
	Update(ctx context.Context, userID, taskID uuid.UUID, input *model.TaskInput) (*model.Task, error)

	// Delete — мягкое удаление (deleted_at). userID должен быть участником команды задачи.
	Delete(ctx context.Context, userID, taskID uuid.UUID) error

	// Restore снимает пометку удаления. userID должен быть участником команды задачи.
	Restore(ctx context.Context, userID, taskID uuid.UUID) (*model.Task, error)

	// GetHistory возвращает историю изменений задачи. userID должен быть участником команды задачи.
	GetHistory(ctx context.Context, taskID, userID uuid.UUID) ([]*model.TaskHistory, error)
}

// TaskReportService — отчёты по задачам и командам (статистика, топ создателей, некорректные assignee).
// Данные через ReportRepository (в модуле task); результаты фильтруются по командам, где userID — участник.
type TaskReportService interface {
	// TeamTaskStats — по каждой команде userID: название, кол-во участников, кол-во задач done за since.
	TeamTaskStats(ctx context.Context, userID uuid.UUID, since time.Time) ([]*model.TeamTaskStats, error)

	// TopCreatorsByTeam — топ-N по созданным задачам в каждой команде за период (только команды userID).
	TopCreatorsByTeam(ctx context.Context, userID uuid.UUID, since time.Time, limit int) ([]*model.TeamTopCreator, error)

	// TasksWithInvalidAssignee — задачи с assignee не из команды (только по командам userID).
	TasksWithInvalidAssignee(ctx context.Context, userID uuid.UUID) ([]*model.Task, error)
}
