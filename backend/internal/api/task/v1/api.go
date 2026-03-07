package task_v1

import (
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/service"
)

// API — HTTP-слой задач: создание, список, получение, обновление, история, отчёты, комментарии.
// Требует JWT (user_id из контекста).
type API struct {
	taskService    service.TaskService
	reportService  service.TaskReportService
	commentService service.TaskCommentService
}

// NewAPI создаёт API.
func NewAPI(
	taskService service.TaskService,
	reportService service.TaskReportService,
	commentService service.TaskCommentService,
) *API {
	return &API{
		taskService:    taskService,
		reportService:  reportService,
		commentService: commentService,
	}
}
