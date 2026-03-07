package report

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

func (s *taskReportService) TasksWithInvalidAssignee(ctx context.Context, userID uuid.UUID) ([]*model.Task, error) {
	userTeams, err := s.teamSvc.ListByUserID(ctx, userID)
	if err != nil {
		logger.Error(ctx, "ListByUserID failed", zap.Error(err))
		return nil, err
	}
	teamIDs := make(map[uuid.UUID]struct{})
	for _, t := range userTeams {
		teamIDs[t.ID] = struct{}{}
	}

	all, err := s.reportRepo.TasksWithInvalidAssignee(ctx, nil)
	if err != nil {
		logger.Error(ctx, "TasksWithInvalidAssignee failed", zap.Error(err))
		return nil, err
	}

	out := make([]*model.Task, 0, len(all))
	for _, task := range all {
		if _, ok := teamIDs[task.TeamID]; ok {
			out = append(out, task)
		}
	}
	return out, nil
}
