package report

import (
	taskRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository"
	svc "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/service"
	teamServiceDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/service"
)

var _ svc.TaskReportService = (*taskReportService)(nil)

// taskReportService — сервис отчётов по задачам.
type taskReportService struct {
	reportRepo taskRepo.ReportRepository
	teamSvc    teamServiceDef.TeamService
}

// NewTaskReportService создаёт сервис отчётов. teamSvc нужен для списка команд пользователя (ListByUserID).
func NewTaskReportService(
	reportRepo taskRepo.ReportRepository,
	teamSvc teamServiceDef.TeamService,
) svc.TaskReportService {
	return &taskReportService{
		reportRepo: reportRepo,
		teamSvc:    teamSvc,
	}
}
