package di

import (
	"context"
	"fmt"

	task_v1 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/task/v1"
	taskRepoDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository"
	taskRepoComment "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/adapter/comment"
	taskRepoHistory "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/adapter/history"
	taskRepoReport "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/adapter/report"
	taskRepoTask "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/adapter/task"
	taskRepoListCache "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/cache/list"
	taskRepoCommentReader "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/repository/comment/reader"
	taskRepoCommentWriter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/repository/comment/writer"
	taskRepoHistoryReader "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/repository/history/reader"
	taskRepoHistoryWriter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/repository/history/writer"
	taskRepoReportReader "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/repository/report"
	taskRepoTaskReader "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/repository/task/reader"
	taskRepoTaskWriter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/repository/task/writer"
	taskServiceDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/service"
	taskServiceComment "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/service/comment"
	taskServiceReport "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/service/report"
	taskServiceImpl "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/service/task"
)

// TaskV1API возвращает HTTP API task v1 (tasks, reports, comments).
func (d *Container) TaskV1API(ctx context.Context) (*task_v1.API, error) {
	if d.taskAPI != nil {
		return d.taskAPI, nil
	}
	taskSvc, err := d.TaskService(ctx)
	if err != nil {
		return nil, fmt.Errorf("task service: %w", err)
	}
	reportSvc, err := d.TaskReportService(ctx)
	if err != nil {
		return nil, fmt.Errorf("task report service: %w", err)
	}
	commentSvc, err := d.TaskCommentService(ctx)
	if err != nil {
		return nil, fmt.Errorf("task comment service: %w", err)
	}
	d.taskAPI = task_v1.NewAPI(taskSvc, reportSvc, commentSvc)
	return d.taskAPI, nil
}

// TaskService возвращает сервис задач.
func (d *Container) TaskService(ctx context.Context) (taskServiceDef.TaskService, error) {
	if d.taskService != nil {
		return d.taskService, nil
	}
	taskRepo, err := d.TaskRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("task repo: %w", err)
	}
	historyRepo, err := d.TaskHistoryRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("task history repo: %w", err)
	}
	memberRepo, err := d.MemberRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("member repo: %w", err)
	}
	txMgr, err := d.UserTxManager(ctx)
	if err != nil {
		return nil, fmt.Errorf("tx manager: %w", err)
	}
	d.taskService = taskServiceImpl.NewTaskService(taskRepo, historyRepo, memberRepo, txMgr)
	return d.taskService, nil
}

// TaskReportService возвращает сервис отчётов по задачам.
func (d *Container) TaskReportService(ctx context.Context) (taskServiceDef.TaskReportService, error) {
	if d.taskReportService != nil {
		return d.taskReportService, nil
	}
	reportRepo, err := d.TaskReportRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("task report repo: %w", err)
	}
	teamSvc, err := d.TeamService(ctx)
	if err != nil {
		return nil, fmt.Errorf("team service: %w", err)
	}
	d.taskReportService = taskServiceReport.NewTaskReportService(reportRepo, teamSvc)
	return d.taskReportService, nil
}

// TaskCommentService возвращает сервис комментариев к задачам.
func (d *Container) TaskCommentService(ctx context.Context) (taskServiceDef.TaskCommentService, error) {
	if d.taskCommentService != nil {
		return d.taskCommentService, nil
	}
	taskRepo, err := d.TaskRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("task repo: %w", err)
	}
	commentRepo, err := d.TaskCommentRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("task comment repo: %w", err)
	}
	memberRepo, err := d.MemberRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("member repo: %w", err)
	}
	txMgr, err := d.UserTxManager(ctx)
	if err != nil {
		return nil, fmt.Errorf("tx manager: %w", err)
	}
	d.taskCommentService = taskServiceComment.NewCommentService(taskRepo, commentRepo, memberRepo, txMgr)
	return d.taskCommentService, nil
}

// TaskRepository возвращает репозиторий задач (адаптер reader + writer + кеш списка).
func (d *Container) TaskRepository(ctx context.Context) (taskRepoDef.TaskRepository, error) {
	if d.taskRepo != nil {
		return d.taskRepo, nil
	}
	db, err := d.SqlxDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("sqlx db: %w", err)
	}
	reader := taskRepoTaskReader.NewRepository(db)
	writer := taskRepoTaskWriter.NewRepository(db)
	var listCache taskRepoDef.TaskListCacheRepository
	if redisClient, err := d.RedisClient(ctx); err == nil {
		listCache = taskRepoListCache.NewRepository(redisClient)
	}
	d.taskRepo = taskRepoTask.NewAdapter(reader, writer, listCache)
	return d.taskRepo, nil
}

// TaskHistoryRepository возвращает репозиторий истории задач.
func (d *Container) TaskHistoryRepository(ctx context.Context) (taskRepoDef.TaskHistoryRepository, error) {
	if d.taskHistoryRepo != nil {
		return d.taskHistoryRepo, nil
	}
	db, err := d.SqlxDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("sqlx db: %w", err)
	}
	reader := taskRepoHistoryReader.NewRepository(db)
	writer := taskRepoHistoryWriter.NewRepository(db)
	d.taskHistoryRepo = taskRepoHistory.NewAdapter(reader, writer)
	return d.taskHistoryRepo, nil
}

// TaskCommentRepository возвращает репозиторий комментариев к задачам.
func (d *Container) TaskCommentRepository(ctx context.Context) (taskRepoDef.TaskCommentRepository, error) {
	if d.taskCommentRepo != nil {
		return d.taskCommentRepo, nil
	}
	db, err := d.SqlxDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("sqlx db: %w", err)
	}
	reader := taskRepoCommentReader.NewRepository(db)
	writer := taskRepoCommentWriter.NewRepository(db)
	d.taskCommentRepo = taskRepoComment.NewAdapter(reader, writer)
	return d.taskCommentRepo, nil
}

// TaskReportRepository возвращает репозиторий отчётов по задачам.
func (d *Container) TaskReportRepository(ctx context.Context) (taskRepoDef.ReportRepository, error) {
	if d.taskReportRepo != nil {
		return d.taskReportRepo, nil
	}
	db, err := d.SqlxDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("sqlx db: %w", err)
	}
	reader := taskRepoReportReader.NewRepository(db)
	d.taskReportRepo = taskRepoReport.NewAdapter(reader)
	return d.taskReportRepo, nil
}
